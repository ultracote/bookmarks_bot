package main

import (
	"bookmarks/db/entity"
	bookmarksTag "bookmarks/handler/bookmarks/tag"
	bookmarksUntag "bookmarks/handler/bookmarks/untag"
	"bookmarks/handler/errors"
	"bookmarks/handler/fallback"
	"bookmarks/handler/list"
	"bookmarks/handler/redirect"
	"bookmarks/handler/save"
	"bookmarks/handler/start"
	"bookmarks/handler/tag"
	tagCreate "bookmarks/handler/tag/create"
	tagCreateDone "bookmarks/handler/tag/create/done"
	tagCreateWait "bookmarks/handler/tag/create/wait"
	"bookmarks/handler/tags"
	"bookmarks/routes"
	"bookmarks/service/router"
	"github.com/mymmrac/telego"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	db := panicOnErr2(connectDB())
	panicOnErr(db.AutoMigrate(&entity.Bookmark{}))
	panicOnErr(db.AutoMigrate(&entity.Tag{}))
	panicOnErr(db.AutoMigrate(&entity.BookmarkTag{}))
	panicOnErr(db.AutoMigrate(&entity.UserState{}))

	token := os.Getenv("BOT_TOKEN")
	bot, err := telego.NewBot(token)
	if err != nil {
		panic(err)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	defer bot.StopLongPolling()

	fallbackHandler := fallback.NewHandler()
	errorsHandler := errors.NewHandler()

	routerService := router.NewRouter(db, fallbackHandler, errorsHandler)
	routerService.AddCommandHandler(routes.CmdStart, start.NewHandler())
	routerService.AddCommandHandler(routes.CmdStop,
		router.Wrap(func(request *router.Request) (*router.Response, error) {
			return router.NewResponse().AddMessage(&telego.SendMessageParams{
				Text: "Goodbye!",
			}), nil
		}),
	)
	routerService.AddCommandHandler(routes.CmdTags, tags.NewHandler(db))
	routerService.AddCommandHandler(routes.CmdSave, save.NewHandler(db))
	routerService.AddCommandHandler(routes.CmdList, list.NewHandler(db))

	routerService.AddCallbackHandler(routes.BookmarkTag, bookmarksTag.NewHandler(db))
	routerService.AddCallbackHandler(routes.BookmarkUntag, bookmarksUntag.NewHandler(db))
	routerService.AddCallbackHandler(routes.TagShow, tag.NewHandler(db))
	routerService.AddCallbackHandler(routes.TagCreateWait, tagCreateWait.NewHandler(db))
	routerService.AddCallbackHandler(routes.TagCreateDone, tagCreateDone.NewHandler())
	routerService.AddCallbackHandler(routes.TagCreate, tagCreate.NewHandler(db))

	sender := router.NewSender(bot, routerService)
	redirectHandler := redirect.NewHandler(sender)

	// Loop through all updates when they came
	for update := range updates {
		var chatID int64
		if update.Message != nil {
			chatID = update.Message.Chat.ID
		} else if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
		}
		if chatID == 0 {
			log.Println("Error: chatID is 0")
			continue
		}

		updateCopy := update
		request := router.NewRequest(&updateCopy, chatID)

		_, err := redirectHandler.Handle(request)
		if err != nil {
			log.Println("error on sender.Handle()", err)
		}

		continue
	}
}

func connectDB() (*gorm.DB, error) {
	// Open database connection
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, err
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func panicOnErr2[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
