# Topic_Api v.0.1
___
**Topic_Api** - сервис создания топиков.

**Для корректной работы нам потребуются следующие библиотеки**
1. gorm.
```go
go get -u gorm.io/gorm
```
2. swagger.
```go
go install github.com/swaggo/swag/cmd/swag@latest
```
3. gin.
```go
go get -u github.com/gin-gonic/gin
```
4. jwt.
5. uuid.

## Gorm.
***Для подключения к базе данных*** - используется функционал библиотеки gorm. Код подключения приведен ниже.
```go
var db *gorm.DB

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		DbUser,
        DbPassword,
        DbHost,
        DbPort,
        DbName,
	)

	// подключение к бд
	db, err := gorm.Open("postgres", connectionString)

	DB = db
```
Реализация ***Миграции*** с помощью библиотеки *gorm*

****Пример****:
```go
m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
	    {
			ID: userID.String(),
			// передаем структуру на основании которой будет создана таблица
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&domain.User{}).Error
				if err != nil {
					return err
				}
				return nil
			},
			// это метод отмены миграции ни разу не использовал
			Rollback: func(tx *gorm.DB) error {
				err := tx.DropTable("users").Error
				if err != nil {
					return err
				}
				return nil
			},
		}
        })
```
## Gin.
Создание ***End point*** с помощью библиотеки джин производится следующим образом:
```go
func SetupRouter() *gin.Engine {
    // создание группы маршрутов
    r := gin.New()
    
    // использование созданного или встроенного middleware в gin производится следующим образом
    r.Use(HandlerFunc())
    
    // если вы хотите поделить запроссы на группы для использования разных middleware вы можете воспользоваться:
    ruse := r.Group("/")
    ruse.Use(AnotheHandlerFunc())
    
    // создание end point
    r.POST("relativePath", ...HandlerFunc)
    
    // для инициализации ручек необходимо их прокинуть в main для этого нам необходимо вернуть r тк он является всей группой end point ов
	return r
}
```
***Запуск сервера*** :
```go
// инициализируем end points
r := routes.SetupRouter()

// запускаем сервер с указанными host и port
if err = r.Run(Host + ":" + Port); err != nil {
		log.WithField("component", "run").Panic(err)
	}
}
```

# Функционал:
+ ниже приведен функционал ***Topic_Api*** для взаимодействия с топиками:
  + Создание.
  
  >За **создание топика** отвечает ниже приведенный хэндлер, ему необходимо передать MultipartForm c ключом data который принимает в себя json модели Topic и файл (опционально) который будет логотипом топика. ***функционал админа***
    ![](![img_2.png](img_2.png))
  ```go
  // шаблон json запроса
    type Topic struct {
        Name    string `json:"name"`
        About   string `json:"about"`
        Creator string `json:"creator"`
    }
  
  authRequiredAdmin.POST("/api/topics", topic.CreateTopic)
    ```
     >**Создание Пользователя**: с занесением в базу данных захэшированного пароля.
    ```go
    r.POST("/api/users/register", user.Register)
    ```
     >**Создание сообщения**: необходимо передать id в uri запроса заместо (:id)
    ```go
  authRequired.POST("/api/messages/topic/:id", massages.Post)
  ```
  + Редактирование.
  
  >**Редактирование Топика**: Аналогичен прошлому методу за ислючением того что мы передаем так же и id топика который хотим изменить, передача происходит в query params.  ***функционал админа***
   ```go
   authRequiredAdmin.PUT("/api/topics/:id", topic.ChangeTopic)
   ```
  > **Редактирование Статуса**: позволяет админу изменить статус пользователя с не подтвержденного на подтвержденный.***функционал админа***
    ```go
    authRequiredAdmin.PUT("/api/users/status", user.ChangeStatus)
  ```  
   > **Редактирование Прав**: позволяет админу изменить права пользователя в плоть до выдачи прав администратора.***функционал админа***
    ```go
    authRequiredAdmin.PUT("api/users/permissions", user.SetPerm)
  ```  
  + Удаление.
  
  > **Удаление топика**: по его id которое передается в query params. ***функционал админа***
   ```go
  authRequiredAdmin.DELETE("/api/topics/:id", topic.DeleteTopic)
  ```
  >**Удаление пользователя**: позволяет удалить пользователя по id ***функционал админа***
  ```go
  authRequiredAdmin.DELETE("/api/topics/:id/user/:user_id", topic.DeleteUser)
  ```
  >**Удаление сообщения**: удаление сообщения по его id передача производится в uri запроса.  ***функционал админа***
    ```go
    authRequiredAdmin.DELETE("/api/messages/users_message/:id", massages.Delete)
    ```
  >**Удаление сообщения пользователем**: функционал идентичен за исключением того, что пользователь может удалить только свои сообщения.
    ```go
    authRequired.DELETE("/api/messages/:id", massages.Delete)
    ```
  + Получение
  >**Получение всех пользователей** : при запросе необходимо передать в query params ключи page и limit которые отвечают за страницу и количество элементов отображаемое на этой странице. ***функционал админа***
     ```go
    authRequiredAdmin.GET("/api/users/all", user.GetUsers)
    ```
  >**Получение по ID** : ***функционал админа***
     ```go
    authRequiredAdmin.GET("/api/users/:id", user.GetUser)
    ``` 
  >**Получение сообщений из топика** : в uri необходимо передать (id) топика и если пользователь подтвержден и находится в этом топике ему отобразится массив сообщений.
     ```go
    authRequired.GET("/api/messages/topic/:id", topic.TopicMassages)
    ```
  >**Получение информации о себе** 
     ```go
    authRequired.GET("/api/users", user.GetUserSelf)
    ```
  >**Получение информации о всех топиках** 
    
![](https://cdn.discordapp.com/attachments/587227868145385502/1221882758104748172/image.png?ex=661431ff&is=6601bcff&hm=a8ae4dbad90b80c88aa82f34696f8053ed6c1a31be27a4b358ceba18cdbe1969&)
     
```go
    authRequired.GET("/api/topics", topic.GetTopics)
   ```
  >**Получение информации топике** 
     ```go
    authRequired.GET("/api/topics/:id", topic.GetTopic)
    ```
  >**Подключение к топику** : необходимо передать в uri запроса id топика к которому хочет подключится пользователь. 
     ```go
    authRequired.GET("/api/topics/join/:id", topic.JoinTopic)
    ```
  
