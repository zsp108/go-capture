package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocapture/go-capture/internal/app"
	"github.com/gocapture/go-capture/internal/ui"
)

const (
	Version   = "v0.1.0"
	AppName   = "GoCapture"
	BuildDate = "2026-08-17"
)

func main() {
	showVersion := flag.Bool("v", false, "显示版本信息")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s %s (Build: %s)\n", AppName, Version, BuildDate)
		fmt.Println("纯 Golang 高性能桌面截图软件 (对标 Snipaste / PixPin / CleanShot X)")
		return
	}

	fmt.Println("================================================================")
	fmt.Printf("🚀 %s %s 原生桌面截图软件已启动\n", AppName, Version)
	fmt.Println("================================================================")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 初始化核心状态机与配置
	cfg := app.DefaultConfig()
	coreApp := app.NewApp(cfg)

	// 2. 初始化原生无边框置顶遮罩控制器
	overlay := ui.NewOverlayController()

	// 3. 注册并启动全局热键监听 (F1 / Ctrl+Shift+A 全局截屏, F3 贴图)
	if err := coreApp.InitHotkeys(ctx); err != nil {
		fmt.Printf("[警告] 全局热键初始化提示: %v\n", err)
	} else {
		fmt.Printf("⚡ 全局截屏热键已就绪: %s | 贴图热键: %s\n", cfg.HotkeyCapture, cfg.HotkeyPin)
	}

	// 4. 监听中断退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("\n✅ GoCapture 已在后台常驻运行 (纯原生架构，无 Web 依赖)。")
	fmt.Println("👉 按下全局截屏快捷键或发送中断信号可操作。按 Ctrl+C 退出。")

	<-sigChan

	fmt.Println("\n🛑 正在退出 GoCapture...")
	overlay.Close()
	coreApp.Shutdown()
	time.Sleep(50 * time.Millisecond)
	fmt.Println("👋 已完全退出。")
}
