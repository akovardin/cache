// Пакет cache предоставляет api для работы с кешем
// В данной версии необходимо подумать, что будет при
// одновременном обращении к Data
package server

import (
    "bufio"
    // "errors"
    "4gophers.com/cache/safemap"
    "fmt"
    "io"
    "log"
    "net"
    "strconv"
    "strings"
    // "time"
)

// Имена наших команд
const (
    SET = "set"
    GET = "get"
)

// Константы для ошибок
const (
    CLIENT_ERROR = "CLIENT_ERROR"
)

// Структура для для хранения данных кеша и метаинформации
type Item struct {
    Key     string
    Flags   int32
    Exptime int
    Length  int
    Data    []byte
}

// Нам нужен мап с конкурентным доступом
var storage safemap.SafeMap

// Инициализируем память для нашего хеша
func init() {
    storage = safemap.New()
}

type Command interface {
    Run()
}

// ParseTextCommand метод который парсит строку с параметрами комманды комманды
func ParseTextCommand(line string, command Command, parse func() error) error {
    log.Print("строка команды ", line)
    // _, err := fmt.Sscanf(line, "set %s %d %d %d\n", &s.Key, &s.Flags, &s.Exptime, &s.Length)

    err := parse()

    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

// Структура для команды set
type SetCommand struct {
    Key     string
    Flags   int32
    Exptime int
    Length  int
    Text    string
    Conn    net.Conn
}

// Записываем в кеш
func (s *SetCommand) Run() {
    log.Println("запустили комманду set")
    // разбираем параметры комманды
    err := ParseTextCommand(s.Text, s, func() error {
        // _, err := fmt.Sscanf(s.Text, "set %s %d %d %d\n", &s.Key, &s.Flags, &s.Exptime, &s.Length)
        _, err := fmt.Sscanf(s.Text, "set %s %d\r\n", &s.Key, &s.Length)
        return err
    })

    log.Println("команда распаршена, ключ:", s.Key)

    // потом получаем данные
    if err == nil {
        reader := bufio.NewReader(s.Conn)
        data, err := reader.Peek(s.Length)

        // Проверяем, не слишком ли короткое сообщение
        if err != nil {
            s.Conn.Write([]byte(CLIENT_ERROR + "\n"))
            return
        }

        // Проверяем, не слишком ли длинное сообщение
        control, err := reader.Peek(s.Length + 2)
        if err != nil {
            s.Conn.Write([]byte(CLIENT_ERROR + "\n"))
            return
        }

        log.Println(data)
        log.Println(control)

        if !strings.HasSuffix(string(control), "\r\n") {
            s.Conn.Write([]byte(CLIENT_ERROR + "\n"))
            return
        }

        storage.Insert(s.Key, Item{Key: s.Key, Length: s.Length, Data: data})
        s.Conn.Write([]byte("STORED\n"))
    }
}

// GetCommand структура с информацией для получения кеша из памяти
type GetCommand struct {
    Name string
    Key  string
    Text string
    Conn net.Conn
}

// Получаем данные из нашего кеша
func (g *GetCommand) Run() {
    log.Println("запустили комманду get")
    // разбираем параметры комманды
    err := ParseTextCommand(g.Text, g, func() error {
        _, err := fmt.Sscanf(g.Text, "get %s\n", &g.Key)
        return err
    })

    if err == nil {
        data, ok := storage.Find(g.Key)
        item := data.(Item)
        if ok {
            // Необходимо для адекватного переноса, так как при считывании
            // последний перенос не учитывался
            g.Conn.Write([]byte("VALUE " + g.Key + " " + strconv.Itoa(item.Length) + "\r\n"))
            g.Conn.Write(item.Data)
            g.Conn.Write([]byte("\r\n"))
            g.Conn.Write([]byte("END\r\n"))
        }
    }
}

// ConnectionHandler обработчик соединения
func ConnectionHandler(conn net.Conn) {
    for {
        // Сначала получаем команду
        command, err := bufio.NewReader(conn).ReadString('\n')

        if err != nil {
            if err == io.EOF {
                log.Println("Error io.EOF", err)
                // Конец сообщения, выходим из цикла получения сообщения
                break
            } else {
                // Не конец сообщения, просто ошибка
                log.Println("Error reading:", err)
            }
        }

        // проверяем, не является ли это командой set
        if strings.HasPrefix(command, SET) {
            // set - указатель на SetCommand. У этого указателя есть метод
            // Run. А значит, что он соответствует интерфейсу Command.
            // Если бы мы создали set := SetCommand{}, то интерфейс был бы другой,
            // так как метод привязан к указателю а не к занчению
            set := &SetCommand{
                Text: command,
                Conn: conn,
            }

            set.Run()
        }

        // проверяем, не является ли это командой get
        if strings.HasPrefix(command, GET) {
            get := &GetCommand{
                Text: command,
                Conn: conn,
            }

            get.Run()
        }
    }
}
