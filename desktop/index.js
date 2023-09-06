const { app, BrowserWindow, Menu } = require('electron');

app.on('ready', () => {

  const loadingWindow = new BrowserWindow({ width: 400, height: 300, frame: false });
  const mainWindow = new BrowserWindow({
    width: 1,
    height: 1,
  });
  
  // 先隐藏主窗口
  mainWindow.hide()
  // 加载第三方网站
  mainWindow.loadURL('https://ai.r9it.com');
  // 加载 loading.html 文件
  loadingWindow.loadFile('loading.html');
  // 隐藏菜单
  Menu.setApplicationMenu(null);
  
  // 监听 loading.html 窗口的 'show-main-window' 事件
  mainWindow.webContents.on('did-finish-load', () => {
	// 最大化窗口
	mainWindow.maximize();
    // 显示主窗口
	mainWindow.show();
	// 关闭加载窗口
    loadingWindow.close();
  });
});
