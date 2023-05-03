package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func db_page_save(strFile string) {
	// 连接SQLite数据库
	db, err := gorm.Open(sqlite.Open(strFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying SQL database: %v", err)
	}
	defer sqlDB.Close()

	// 自动迁移数据库结构
	err = db.AutoMigrate(&PageMainData{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	db.Begin()
	for _, item := range g_PageMainData {
		db.Create(item)
	}
	db.Commit()
}

func db_page_load(strFile string) {
	// 连接SQLite数据库
	db, err := gorm.Open(sqlite.Open(strFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying SQL database: %v", err)
	}
	defer sqlDB.Close()

	// 自动迁移数据库结构
	err = db.AutoMigrate(&PageMainData{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = db.Order("draw_date").Find(&g_PageMainData).Error
}
