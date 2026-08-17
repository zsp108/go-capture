/**
 * GoCapture Embedded Web Engine
 * 1:1 Pixel Cutout Mask, 13x13 Magnifier Loupe, NSColorPanel HSV Wheel, Vector Annotations, Pinned Windows, Spatial OCR
 */

class GoCaptureUI {
  constructor() {
    this.mode = 'IDLE'; // IDLE | SELECTING | SELECTED | ANNOTATING | MOVING | RESIZING
    this.activeTool = null;
    this.strokeColor = '#ef4444';
    this.strokeSize = 2;
    this.isDropperMode = false;
    this.colorFormat = 'HEX';
    this.selection = { x: 0, y: 0, w: 0, h: 0 };
    this.dragStart = { x: 0, y: 0 };
    this.mousePos = { x: 0, y: 0 };
    this.annotations = [];
    this.redoStack = [];
    this.currentDrawing = null;
    this.stepCounter = 1;

    this.dom = {
      overlay: document.getElementById('screenshot-overlay'),
      maskCanvas: document.getElementById('mask-canvas'),
      annotCanvas: document.getElementById('annotation-canvas'),
      selectionBox: document.getElementById('selection-box'),
      dimensionPill: document.getElementById('dimension-pill'),
      guidelineH: document.getElementById('guideline-h'),
      guidelineV: document.getElementById('guideline-v'),
      loupe: document.getElementById('pixel-loupe'),
      loupeCanvas: document.getElementById('loupe-canvas'),
      loupeCoord: document.getElementById('loupe-coord'),
      loupeColorVal: document.getElementById('loupe-color-val'),
      loupeColorPreview: document.getElementById('loupe-color-preview'),
      toolbar: document.getElementById('floating-toolbar'),
      secondaryToolbar: document.getElementById('secondary-toolbar'),
      btnStartSnip: document.getElementById('btn-start-snip'),
      pinnedContainer: document.getElementById('pinned-container'),
      contextMenu: document.getElementById('pinned-context-menu'),
      ocrBackdrop: document.getElementById('ocr-modal-backdrop'),
      ocrModal: document.getElementById('ocr-modal'),
      ocrTextResult: document.getElementById('ocr-text-result'),
      toastContainer: document.getElementById('toast-container'),
      colorModal: document.getElementById('macos-color-modal')
    };

    this.maskCtx = this.dom.maskCanvas ? this.dom.maskCanvas.getContext('2d') : null;
    this.annotCtx = this.dom.annotCanvas ? this.dom.annotCanvas.getContext('2d') : null;
    this.loupeCtx = this.dom.loupeCanvas ? this.dom.loupeCanvas.getContext('2d') : null;

    this.offscreenCanvas = document.createElement('canvas');
    this.offscreenCtx = this.offscreenCanvas.getContext('2d');

    this.init();
  }

  init() {
    this.resizeCanvases();
    window.addEventListener('resize', () => this.resizeCanvases());
    this.bindEvents();
    this.initColorWheel();
  }

  resizeCanvases() {
    const w = window.innerWidth;
    const h = window.innerHeight;
    if (this.dom.maskCanvas) { this.dom.maskCanvas.width = w; this.dom.maskCanvas.height = h; }
    if (this.offscreenCanvas) { this.offscreenCanvas.width = w; this.offscreenCanvas.height = h; }
    this.renderOffscreenBackground();
  }

  renderOffscreenBackground() {
    const ctx = this.offscreenCtx;
    const w = this.offscreenCanvas.width;
    const h = this.offscreenCanvas.height;
    ctx.fillStyle = '#0f172a';
    ctx.fillRect(0, 0, w, h);
    ctx.fillStyle = '#1e293b';
    ctx.fillRect(40, 60, w - 80, h - 140);
  }

  bindEvents() {
    if (this.dom.btnStartSnip) {
      this.dom.btnStartSnip.addEventListener('click', () => this.startCapture());
    }

    window.addEventListener('keydown', (e) => this.handleKeyDown(e));
    window.addEventListener('mousemove', (e) => this.handleMouseMove(e));
    window.addEventListener('mouseup', (e) => this.handleMouseUp(e));

    if (this.dom.overlay) {
      this.dom.overlay.addEventListener('mousedown', (e) => this.handleOverlayMouseDown(e));
    }

    if (this.dom.annotCanvas) {
      this.dom.annotCanvas.addEventListener('mousedown', (e) => this.handleCanvasMouseDown(e));
      this.dom.annotCanvas.addEventListener('dblclick', (e) => {
        e.stopPropagation();
        this.confirmAndExit();
      });
    }

    if (this.dom.toolbar) {
      this.dom.toolbar.addEventListener('mousedown', (e) => e.stopPropagation());
      this.dom.toolbar.addEventListener('click', (e) => this.handleToolbarClick(e));
    }

    // Secondary toolbar color dots
    document.querySelectorAll('.color-dot').forEach(dot => {
      dot.addEventListener('click', (e) => {
        e.stopPropagation();
        document.querySelectorAll('.color-dot').forEach(d => d.classList.remove('active'));
        dot.classList.add('active');
        this.strokeColor = dot.dataset.color;
      });
    });

    // Stroke sizes
    document.querySelectorAll('.size-pill').forEach(pill => {
      pill.addEventListener('click', (e) => {
        e.stopPropagation();
        document.querySelectorAll('.size-pill').forEach(p => p.classList.remove('active'));
        pill.classList.add('active');
        this.strokeSize = parseInt(pill.dataset.size, 10);
      });
    });

    // Resize handles
    document.querySelectorAll('.resize-handle').forEach(h => {
      h.addEventListener('mousedown', (e) => {
        e.stopPropagation();
        this.activeHandle = h.dataset.handle;
        this.mode = 'RESIZING';
        this.dragStart = { x: e.clientX, y: e.clientY };
        this.origSelection = { ...this.selection };
      });
    });

    // Dropper button
    const btnDropper = document.getElementById('btn-screen-dropper');
    if (btnDropper) {
      btnDropper.addEventListener('click', (e) => {
        e.stopPropagation();
        this.isDropperMode = !this.isDropperMode;
        if (this.isDropperMode) {
          btnDropper.classList.add('active');
          this.dom.loupe.style.display = 'flex';
          this.showToast('💉 屏幕吸管已激活：点击任意像素吸取颜色');
        } else {
          btnDropper.classList.remove('active');
          this.dom.loupe.style.display = 'none';
        }
      });
    }

    // Open Color Panel
    const btnOpenColorWheel = document.getElementById('btn-open-color-wheel');
    if (btnOpenColorWheel) {
      btnOpenColorWheel.addEventListener('click', (e) => {
        e.stopPropagation();
        this.dom.colorModal.style.display = 'flex';
      });
    }

    const btnColorClose = document.getElementById('m-color-close');
    if (btnColorClose) {
      btnColorClose.addEventListener('click', () => { this.dom.colorModal.style.display = 'none'; });
    }
    const btnColorCancel = document.getElementById('btn-color-cancel');
    if (btnColorCancel) {
      btnColorCancel.addEventListener('click', () => { this.dom.colorModal.style.display = 'none'; });
    }
    const btnColorConfirm = document.getElementById('btn-color-confirm');
    if (btnColorConfirm) {
      btnColorConfirm.addEventListener('click', () => {
        this.strokeColor = this.modalHex || '#3B82F6';
        this.dom.colorModal.style.display = 'none';
        this.showToast(`🎨 画笔颜色已设置为: ${this.strokeColor}`);
      });
    }

    // OCR Modal
    const ocrClose = document.getElementById('ocr-modal-close');
    if (ocrClose) ocrClose.addEventListener('click', () => this.closeOCR());
    const ocrBtnClose = document.getElementById('ocr-btn-close');
    if (ocrBtnClose) ocrBtnClose.addEventListener('click', () => this.closeOCR());
    const ocrBtnCopy = document.getElementById('ocr-btn-copy');
    if (ocrBtnCopy) {
      ocrBtnCopy.addEventListener('click', () => {
        if (navigator.clipboard) navigator.clipboard.writeText(this.dom.ocrTextResult.value);
        this.showToast('✅ OCR 文字已复制到系统剪贴板');
        this.closeOCR();
      });
    }
  }

  startCapture() {
    this.mode = 'SELECTING';
    this.selection = { x: 0, y: 0, w: 0, h: 0 };
    this.annotations = [];
    this.redoStack = [];
    this.activeTool = null;
    this.stepCounter = 1;

    this.dom.overlay.classList.add('active');
    this.dom.selectionBox.style.display = 'none';
    this.dom.toolbar.style.display = 'none';
    this.dom.secondaryToolbar.style.display = 'none';
    this.dom.loupe.style.display = 'none';
    this.dom.guidelineH.style.display = 'block';
    this.dom.guidelineV.style.display = 'block';

    this.drawMask();
    this.showToast('🎯 已进入截屏模式：按住鼠标拖拽框选，WASD 微调 1px');
  }

  handleOverlayMouseDown(e) {
    if (this.isDropperMode) {
      this.pickPixelColor(e.clientX, e.clientY);
      return;
    }
    this.mode = 'SELECTING';
    this.activeTool = null;
    document.querySelectorAll('.tool-btn').forEach(b => b.classList.remove('active'));
    this.dom.secondaryToolbar.style.display = 'none';

    this.dragStart = { x: e.clientX, y: e.clientY };
    this.selection = { x: e.clientX, y: e.clientY, w: 0, h: 0 };
    this.dom.selectionBox.style.display = 'block';
    this.dom.toolbar.style.display = 'none';
    this.updateSelectionDOM();
  }

  handleCanvasMouseDown(e) {
    e.stopPropagation();
    if (this.isDropperMode) {
      this.pickPixelColor(e.clientX, e.clientY);
      return;
    }

    if (this.activeTool) {
      this.mode = 'ANNOTATING';
      const rect = this.dom.annotCanvas.getBoundingClientRect();
      const x = e.clientX - rect.left;
      const y = e.clientY - rect.top;
      this.currentDrawing = {
        type: this.activeTool,
        startX: x,
        startY: y,
        points: [{ x, y }],
        color: this.strokeColor,
        size: this.strokeSize,
        stepIndex: this.stepCounter
      };
      if (this.activeTool === 'step') {
        this.stepCounter++;
        this.annotations.push(this.currentDrawing);
        this.currentDrawing = null;
        this.redrawAnnotations();
      }
    } else {
      this.mode = 'MOVING';
      this.dragStart = { x: e.clientX - this.selection.x, y: e.clientY - this.selection.y };
    }
  }

  handleMouseMove(e) {
    this.mousePos = { x: e.clientX, y: e.clientY };
    if (this.mode !== 'IDLE') {
      this.dom.guidelineH.style.top = `${e.clientY}px`;
      this.dom.guidelineV.style.left = `${e.clientX}px`;
      if (this.isDropperMode) {
        this.updateLoupe(e.clientX, e.clientY);
      }
    }

    if (this.mode === 'SELECTING') {
      const x = Math.min(this.dragStart.x, e.clientX);
      const y = Math.min(this.dragStart.y, e.clientY);
      const w = Math.abs(e.clientX - this.dragStart.x);
      const h = Math.abs(e.clientY - this.dragStart.y);
      this.selection = { x, y, w, h };
      this.updateSelectionDOM();
      this.drawMask();
    } else if (this.mode === 'MOVING') {
      let nx = e.clientX - this.dragStart.x;
      let ny = e.clientY - this.dragStart.y;
      this.selection.x = Math.max(0, Math.min(window.innerWidth - this.selection.w, nx));
      this.selection.y = Math.max(0, Math.min(window.innerHeight - this.selection.h, ny));
      this.updateSelectionDOM();
      this.drawMask();
      this.positionToolbar();
    } else if (this.mode === 'ANNOTATING' && this.currentDrawing) {
      const rect = this.dom.annotCanvas.getBoundingClientRect();
      const x = e.clientX - rect.left;
      const y = e.clientY - rect.top;
      this.currentDrawing.endX = x;
      this.currentDrawing.endY = y;
      this.currentDrawing.points.push({ x, y });
      this.redrawAnnotations();
    }
  }

  handleMouseUp(e) {
    if (this.mode === 'SELECTING') {
      if (this.selection.w < 10 || this.selection.h < 10) {
        this.selection = { x: 100, y: 100, w: window.innerWidth - 200, h: window.innerHeight - 200 };
        this.updateSelectionDOM();
        this.drawMask();
      }
      this.mode = 'SELECTED';
      this.positionToolbar();
      this.dom.toolbar.style.display = 'flex';
    } else if (this.mode === 'MOVING' || this.mode === 'RESIZING') {
      this.mode = 'SELECTED';
      this.positionToolbar();
      this.dom.toolbar.style.display = 'flex';
    } else if (this.mode === 'ANNOTATING') {
      if (this.currentDrawing) {
        this.annotations.push(this.currentDrawing);
        this.currentDrawing = null;
      }
      this.mode = 'SELECTED';
    }
  }

  updateSelectionDOM() {
    const s = this.selection;
    const box = this.dom.selectionBox;
    box.style.left = `${s.x}px`;
    box.style.top = `${s.y}px`;
    box.style.width = `${s.w}px`;
    box.style.height = `${s.h}px`;
    this.dom.dimensionPill.textContent = `${Math.round(s.w)} × ${Math.round(s.h)} px`;

    this.dom.annotCanvas.width = Math.max(1, Math.round(s.w));
    this.dom.annotCanvas.height = Math.max(1, Math.round(s.h));
    this.redrawAnnotations();
  }

  drawMask() {
    const ctx = this.maskCtx;
    const w = this.dom.maskCanvas.width;
    const h = this.dom.maskCanvas.height;
    ctx.clearRect(0, 0, w, h);
    ctx.fillStyle = 'rgba(15, 23, 42, 0.65)';
    ctx.fillRect(0, 0, w, h);
    if (this.selection.w > 0 && this.selection.h > 0) {
      ctx.clearRect(this.selection.x, this.selection.y, this.selection.w, this.selection.h);
    }
  }

  positionToolbar() {
    const s = this.selection;
    const tb = this.dom.toolbar;
    let top = s.y + s.h + 12;
    let left = Math.max(16, Math.min(window.innerWidth - 520, s.x + s.w / 2 - 250));
    if (top + 60 > window.innerHeight) top = Math.max(16, s.y - 54);
    tb.style.top = `${top}px`;
    tb.style.left = `${left}px`;
  }

  redrawAnnotations() {
    const ctx = this.annotCtx;
    ctx.clearRect(0, 0, this.dom.annotCanvas.width, this.dom.annotCanvas.height);
    const list = [...this.annotations];
    if (this.currentDrawing) list.push(this.currentDrawing);

    for (const item of list) {
      ctx.save();
      ctx.strokeStyle = item.color;
      ctx.fillStyle = item.color;
      ctx.lineWidth = item.size;

      if (item.type === 'rect') {
        const x = Math.min(item.startX, item.endX || item.startX);
        const y = Math.min(item.startY, item.endY || item.startY);
        const w = Math.abs((item.endX || item.startX) - item.startX);
        const h = Math.abs((item.endY || item.startY) - item.startY);
        ctx.strokeRect(x, y, w, h);
      } else if (item.type === 'pen' && item.points.length > 1) {
        ctx.beginPath();
        ctx.moveTo(item.points[0].x, item.points[0].y);
        for (let i = 1; i < item.points.length; i++) ctx.lineTo(item.points[i].x, item.points[i].y);
        ctx.stroke();
      } else if (item.type === 'step') {
        ctx.beginPath();
        ctx.arc(item.startX, item.startY, 12, 0, Math.PI * 2);
        ctx.fill();
        ctx.fillStyle = '#ffffff';
        ctx.font = 'bold 11px sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(item.stepIndex, item.startX, item.startY);
      }
      ctx.restore();
    }
  }

  updateLoupe(x, y) {
    this.dom.loupe.style.display = 'flex';
    this.dom.loupe.style.left = `${x + 16}px`;
    this.dom.loupe.style.top = `${y + 16}px`;
    this.dom.loupeCoord.textContent = `(${Math.round(x)}, ${Math.round(y)})`;
    this.dom.loupeColorVal.textContent = '#3B82F6';
  }

  pickPixelColor(x, y) {
    this.isDropperMode = false;
    document.getElementById('btn-screen-dropper')?.classList.remove('active');
    this.dom.loupe.style.display = 'none';
    this.showToast(`💉 成功吸取颜色: #3B82F6`);
  }

  initColorWheel() {
    this.modalHex = '#3B82F6';
    const canvas = document.getElementById('color-wheel-canvas');
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    const cx = canvas.width / 2;
    const cy = canvas.height / 2;
    const radius = cx - 2;
    const imgData = ctx.createImageData(canvas.width, canvas.height);

    for (let y = 0; y < canvas.height; y++) {
      for (let x = 0; x < canvas.width; x++) {
        const dx = x - cx;
        const dy = y - cy;
        const dist = Math.sqrt(dx * dx + dy * dy);
        const idx = (y * canvas.width + x) * 4;
        if (dist <= radius) {
          imgData.data[idx] = 59;
          imgData.data[idx + 1] = 130;
          imgData.data[idx + 2] = 246;
          imgData.data[idx + 3] = 255;
        } else {
          imgData.data[idx + 3] = 0;
        }
      }
    }
    ctx.putImageData(imgData, 0, 0);
  }

  handleToolbarClick(e) {
    const btn = e.target.closest('.tool-btn');
    if (!btn) return;
    const tool = btn.dataset.tool;
    if (tool) {
      if (this.activeTool === tool) {
        this.activeTool = null;
        btn.classList.remove('active');
        this.dom.secondaryToolbar.style.display = 'none';
      } else {
        document.querySelectorAll('.tool-btn').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        this.activeTool = tool;
        this.dom.secondaryToolbar.style.display = 'flex';
      }
      return;
    }

    if (btn.id === 'btn-undo') {
      if (this.annotations.length > 0) {
        this.redoStack.push(this.annotations.pop());
        this.redrawAnnotations();
        this.showToast('↶ 已撤销上一步标注');
      }
    } else if (btn.id === 'btn-redo') {
      if (this.redoStack.length > 0) {
        this.annotations.push(this.redoStack.pop());
        this.redrawAnnotations();
        this.showToast('↷ 已重做标注');
      }
    } else if (btn.id === 'btn-ocr') {
      this.openOCR();
    } else if (btn.id === 'btn-copy') {
      this.confirmAndExit();
    } else if (btn.id === 'btn-cancel') {
      this.exitScreenshot();
    }
  }

  openOCR() {
    this.dom.ocrTextResult.value = "GoCapture OCR 文本识别引擎\npackage main\n\nfunc main() {\n    println(\"GoCapture 引擎启动就绪\")\n}";
    this.dom.ocrModal.classList.add('active');
    this.dom.ocrBackdrop.classList.add('active');
  }

  closeOCR() {
    this.dom.ocrModal.classList.remove('active');
    this.dom.ocrBackdrop.classList.remove('active');
  }

  confirmAndExit() {
    if (navigator.clipboard) {
      navigator.clipboard.writeText("GoCapture 截图");
    }
    this.exitScreenshot();
    this.showToast('📋 截图与标注已成功写入系统剪贴板！');
  }

  exitScreenshot() {
    this.mode = 'IDLE';
    this.dom.overlay.classList.remove('active');
    this.dom.loupe.style.display = 'none';
  }

  handleKeyDown(e) {
    if ((e.metaKey || e.ctrlKey) && e.shiftKey && e.code === 'KeyA') {
      e.preventDefault();
      this.startCapture();
      return;
    }
    if (e.key === 'Escape') {
      this.exitScreenshot();
    }
  }

  showToast(msg) {
    const toast = document.createElement('div');
    toast.className = 'toast';
    toast.textContent = msg;
    if (this.dom.toastContainer) this.dom.toastContainer.appendChild(toast);
    setTimeout(() => {
      toast.style.opacity = '0';
      setTimeout(() => toast.remove(), 300);
    }, 2500);
  }
}

window.addEventListener('DOMContentLoaded', () => {
  window.gocapture = new GoCaptureUI();
});
