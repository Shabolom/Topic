package main

import (
	"Arkadiy_Servis_authorization/config"
	_ "Arkadiy_Servis_authorization/docs"
	migrate "Arkadiy_Servis_authorization/init"
	"Arkadiy_Servis_authorization/iternal/routes"
	"Arkadiy_Servis_authorization/iternal/tools"

	log "github.com/sirupsen/logrus"
)

func main() {
	//	@title		User API
	//	@version	1.0.0

	// 	@description 	Это выпускной проэкт с использованием свагера
	// 	@termsOfService  сдесь были бы условия использования еслиб я их мог обозначить
	// 	@contact.url    тут моя контактная информация (https://vk.com/id192672036)
	// 	@contact.email  tima.gorenskiy@mail.ru

	// 	@securityDefinitions.apikey  ApiKeyAuth
	//  @in header
	//  @name Authorization

	//	@host		localhost:8000

	config.CheckFlagEnv()
	tools.InitLogger()

	// config.InitPgSQL инициализируем подключение к базе данных
	err := config.InitPgSQL()
	if err != nil {
		log.WithField("component", "initialization").Panic(err)
	}

	// вызываем миграцию структуры в базу данных
	migrate.Migrate()

	//test.ClientGet()
	//test.Redirect()

	// конфигурация (инициализация) end point или ручка (можно назвать имя запроса)
	// (как api student) URLов пример (localhost, 8080) конфигурация всех URLов которые будет
	// обрабатывать сервер

	r := routes.SetupRouter()

	// запуск сервера
	if err = r.Run(config.Env.Host + ":" + config.Env.Port); err != nil {
		log.WithField("component", "run").Panic(err)
	}

}
