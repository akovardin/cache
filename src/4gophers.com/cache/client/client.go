package client

// Get используется для получения данных и кеша
func Get(key string, result interface{}, executionLimit, expire time.Duration, getter func()) error {
    // getter нужно вызывать только если expire прогорел
    getter()
    return nil
}

func Set(key string, data []byte) {

}
