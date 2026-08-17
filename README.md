# GoCapture (v0.1.0)

> **高性能纯 Golang 原生桌面截图软件**  
> 对标 Snipaste、PixPin 与 CleanShot X，提供毫秒级唤醒、1:1 物理像素无畸变快照、字符级 100% 全包含 OCR 过滤、13×13 苹果级奇数网格取色器与桌面置顶贴图。

---

## ⚡ 核心设计与特性

1. **⚡ 毫秒级极速唤醒 (< 5ms)**：
   - 纯 Go + 原生 CGo 操作系统底层直接抓屏（macOS `CoreGraphics` / Windows `DXGI` / Linux `X11`），0 内存多余拷贝；
2. **🔍 100% 物理像素无畸变快照 (T0 Snapshot Freeze)**：
   - 唤起瞬间冻结屏幕真实物理位图，配合原生透明挖孔遮罩，绝对杜绝字体发虚与拉伸变形；
3. **🎯 字符级微粒度 100% 全包含 OCR 过滤 (Strict 100% Enclosure)**：
   - 计算单个字符物理包围盒，仅当字符 100% 完全处于选区框内部时才予提取，严禁边缘切字与粘连字误识别；
4. **🍎 苹果级 13×13 奇数高精度放大镜 (NSColorSampler)**：
   - 13×13 奇数像素采样矩阵（半宽 6px），中心单元格 `(Col 6, Row 6)` 严丝合缝绝对对齐光标真实物理坐标，按需吸管呼出；
5. **🎨 7 大矢量标注与命令历史栈 (Zero-Interference)**：
   - 矩形、箭头、平滑贝塞尔画笔、荧光笔、真实像素马赛克、自增步骤序号徽标、文字输入；
   - 标注激活后独占输入流，彻底杜绝选区拉伸变形误触；完整支持 Undo (`⌘Z`) / Redo (`⌘Y`)；
6. **📌 桌面置顶多维贴图系统 (Pinned Windows)**：
   - 置顶悬浮，支持鼠标滚轮以 5% 步长在 `20% ~ 300%` 范围内缩放、`10% ~ 100%` 透明度调节、顺时针 90° 旋转、水平镜像翻转及右键菜单；
7. **⌨️ 1 像素键盘微操**：
   - `W / A / S / D`：鼠标十字光标向 上/左/下/右 精确微调 `1px` 并即时刷新取色器；
   - `↑ / ↓ / ← / →`：选区位置向 上/下/左/右 平移 `1px`；
   - `Shift + 方向键`：选区宽高在当前边缘扩展或收缩 `1px`；
   - `⌘/Ctrl + A`：一键全屏覆盖。

---

## 📂 项目模块结构

```
/Users/spz/Downloads/spz/go-capture/
├── go.mod                              # Go 1.22 模块依赖定义
├── Makefile                            # 跨平台构建、测试、打包 .app 指令
├── README.md                           # 详细架构规范、快捷键与使用手册
├── build/
│   └── macos/
│       ├── Info.plist                  # macOS App Bundle 属性配置
│       └── PkgInfo                     # macOS 包标识文件
├── cmd/
│   └── gocapture/
│       └── main.go                     # 原生桌面程序总入口、热键守护与常驻后台
├── pkg/
│   ├── capture/                        # 跨平台 0-拷贝屏幕快照捕获底层 (CGo / OS API)
│   │   ├── capturer.go                 # ScreenCapturer 统一接口 & 图像裁剪工具
│   │   ├── capture_darwin.go           # macOS CoreGraphics (CGDisplayCreateImage) Retina 1:1 无损
│   │   ├── capture_windows.go          # Windows DXGI Desktop Duplication / GDI+ BitBlt
│   │   └── capture_linux.go            # Linux X11 XShmGetImage / Wayland Portal
│   ├── ocr/                        # 字符级微粒度空间 OCR 提取引擎
│   │   ├── models.go                   # CharBoundingBox 字符级包围盒模型
│   │   ├── spatial_filter.go           # 严格 100% 全包含空间包围盒过滤算法
│   │   ├── spatial_filter_test.go      # 单元测试（切字剔除、全包含判定）
│   │   └── ocr_engine.go               # RapidOCR / ONNX 离线推理接口
│   ├── loupe/                      # 苹果级 13×13 奇数网格像素放大镜与取色器
│   │   ├── color_space.go              # HEX / RGB / HSL 色彩空间转换引擎
│   │   ├── color_sampler.go            # 13x13 奇数采样网格（中心点第6行第6列绝对对齐光标）
│   │   └── color_sampler_test.go       # 取色器采样与坐标对齐测试
│   ├── annotation/                 # 7 大矢量标注与命令历史栈
│   │   ├── shape.go                # 矩形、箭头、平滑画笔、荧光笔、马赛克、序号、文字
│   │   ├── command.go              # Command 模式 Undo/Redo 双栈与步骤计数器
│   │   ├── renderer.go             # 矢量图形光栅化合成与 PNG 导出器
│   │   └── annotation_test.go      # 标注历史栈与渲染器测试
│   ├── pin/                        # 桌面置顶多维贴图系统 (Pinned Windows)
│   │   ├── pin_window.go           # 贴图几何变换（缩放20%~300%、透明度10%~100%、旋转、镜像）
│   │   ├── pin_manager.go          # Topmost 置顶窗口生命周期管理器
│   │   └── pin_test.go             # 贴图缩放与几何变换测试
│   ├── clipboard/                  # 系统剪贴板读写引擎 (PNG 位图 / 纯文本)
│   │   └── clipboard.go            # 跨平台剪贴板极速推送
│   └── hotkey/                     # 全局低级键盘热键守护
│       └── hotkey.go               # F1 / Ctrl+Shift+A 全局截屏 / F3 贴图快捷键监听
└── internal/
    ├── app/                            # 核心状态机与业务协调器
    │   ├── app.go                      # 严格状态机 (IDLE, SELECTING, SELECTED, ANNOTATING, PINNED)
    │   ├── config.go                   # 软件配置项
    │   └── events.go                   # 事件总线定义
    └── ui/                             # 桌面端无边框全屏置顶遮罩与原生窗口控制器
        ├── window.go                   # NativeWindow 跨平台原生窗口接口
        ├── window_darwin.go            # macOS Cocoa NSWindow 原生窗口实现
        ├── window_windows.go           # Windows Win32 原生窗口实现
        ├── window_linux.go             # Linux X11 原生窗口实现
        └── overlay.go                  # 原生遮罩协调器
```

---

## 🚀 编译、打包与运行指令

### 1. 编译为 macOS 原生应用程序包 (`.app`)
```bash
make app
# 将在 bin/ 目录下生成标准的 macOS 应用程序：bin/GoCapture.app
```

### 2. 编译并直接启动 `.app`
```bash
make run-app
# 自动编译并调用 macOS 原生 open 命令启动 GoCapture.app
```

### 3. 编译命令行二进制
```bash
make build
# 产物输出于 bin/gocapture
```

### 4. 运行全量单元测试
```bash
make test
# 或
go test -v ./pkg/...
```

---

## ⌨️ 快捷键与操作矩阵

| 操作类型 | 快捷键 | 功能描述 |
| :--- | :--- | :--- |
| **全局截屏** | `F1` / `Ctrl+Shift+A` | 毫秒级抓取屏幕并呼出全屏原生遮罩 |
| **鼠标微操** | `W` / `A` / `S` / `D` | 鼠标光标向 上/左/下/右 微调 1 像素 |
| **选区移动** | `↑` / `↓` / `←` / `→` | 选区位置向 上/下/左/右 平移 1 像素 |
| **选区微调** | `Shift` + 方向键 | 选区边缘向对应方向扩展/收缩 1 像素 |
| **全屏覆盖** | `⌘/Ctrl + A` | 一键最大化选区覆盖当前显示器全屏 |
| **取色复制** | `C` 键 | 复制取色器中心当前像素颜色代码 (HEX/RGB) |
| **格式切换** | `Shift` 键 | 循环切换取色器格式 (HEX / RGB / HSL) |
| **矢量撤销/重做**| `⌘/Ctrl + Z` / `⌘/Ctrl + Y` | 撤销上一步标注 / 恢复重做 |
| **桌面贴图** | `📌` 按钮 / `F3` | 将选区截屏钉在桌面最顶层 |
| **贴图缩放** | 滚轮滑动 (`Wheel`) | 在 `20% ~ 300%` 范围内按 5% 步长缩放贴图 |
| **完成复制** | `Enter` / 选区内双击 | 导出合并标注的 PNG 位图至剪贴板并退出 |
| **退出截屏** | `Esc` | 退出当前截屏遮罩 |

---

*GoCapture - Designed for Precision & Speed | 2026*
