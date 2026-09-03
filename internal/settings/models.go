package settings

// SingletonID — фиксированный id единственной строки smtp_settings (см.
// query/smtp_settings.sql: GetSMTPSettings/UpsertSMTPSettings всегда по
// этому id — синглтон обеспечивается на уровне приложения, не БД).
const SingletonID = "smtp"

// SMTPSettingsResponse — конфиг без пароля, он наружу никогда не отдаётся.
// Configured=false, если настройки ещё не заданы (первый заход) — это не
// ошибка, а нормальное начальное состояние страницы настроек.
type SMTPSettingsResponse struct {
	Configured  bool    `json:"configured"`
	Host        string  `json:"host"`
	Port        int32   `json:"port"`
	Username    string  `json:"username"`
	FromAddress string  `json:"from_address"`
	FromName    *string `json:"from_name"`
	UseTLS      bool    `json:"use_tls"`
}

// UpsertSMTPSettingsRequest — Password опционален: не передан (или пустая
// строка) при уже существующих настройках — прежний пароль сохраняется как
// есть; при первоначальной настройке (настроек ещё нет) обязателен.
type UpsertSMTPSettingsRequest struct {
	Host        string  `json:"host"`
	Port        int32   `json:"port"`
	Username    string  `json:"username"`
	Password    *string `json:"password" description:"Опционален при обновлении — не передан, прежний пароль сохраняется"`
	FromAddress string  `json:"from_address"`
	FromName    *string `json:"from_name"`
	UseTLS      bool    `json:"use_tls" description:"STARTTLS"`
}
