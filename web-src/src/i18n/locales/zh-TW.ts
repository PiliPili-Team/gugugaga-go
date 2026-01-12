export default {
  common: {
    save: '儲存',
    cancel: '取消',
    delete: '刪除',
    add: '新增',
    edit: '編輯',
    confirm: '確認',
    loading: '載入中...',
    seconds: '秒',
    days: '天',
    success: '成功',
    error: '錯誤',
    warning: '警告'
  },
  
  nav: {
    dashboard: '儀表板',
    basic: '基礎設定',
    integrations: '服務整合',
    mappings: '路徑對映',
    ignore: '忽略設定',
    advanced: '進階設定',
    oauth: 'OAuth認證',
    logs: '即時日誌',
    actions: '快捷操作',
    more: '更多'
  },
  
  header: {
    title: 'GD Watcher',
    subtitle: 'Google Drive 檔案監控',
    save: '儲存設定',
    saving: '儲存中...',
    logout: '登出',
    language: '語言',
    theme: '主題色',
    customTheme: '自訂主題',
    appearance: '外觀'
  },
  
  sidebar: {
    collapse: '收合側邊欄',
    expand: '展開側邊欄'
  },
  
  dashboard: {
    systemStatus: '系統狀態',
    online: '運行中',
    offline: '離線',
    uptime: '運行時間',
    activityStream: '活動流',
    historyCompletedTasks: '歷史完成任務',
    todayCompletedTasks: '今日完成任務',
    cpuUsage: 'CPU 負載',
    memoryUsage: '記憶體負載'
  },
  
  login: {
    title: '歡迎回來',
    subtitle: '登入以繼續使用 GD Watcher',
    username: '使用者名稱',
    password: '密碼',
    usernamePlaceholder: '請輸入使用者名稱',
    passwordPlaceholder: '請輸入密碼',
    submit: '登入',
    logging: '登入中...',
    error: '使用者名稱或密碼錯誤',
    footer: 'Google Drive 檔案變更監控系統'
  },
  
  logs: {
    title: '即時日誌',
    autoScroll: '自動捲動',
    clear: '清空日誌',
    empty: '暫無日誌'
  },
  
  actions: {
    title: '快捷操作',
    testPath: '測試路徑',
    testPathPlaceholder: '/雲端硬碟/Test/1.mkv',
    test: '測試',
    sync: '觸發同步',
    syncing: '同步中...',
    rcloneFull: '強刷 Rclone',
    rcloneFulling: '刷新中...',
    rebuildTree: '重建檔案樹',
    rebuildingTree: '重建中...',
    clearPanel: '清空面板',
    clearFiles: '清空日誌檔案',
    confirm: {
      rcloneFull: '確定要強制刷新所有 Rclone 實例嗎？',
      rebuildTree: '確定要重建檔案樹嗎？這可能需要一些時間。',
      clearFiles: '確定要刪除所有日誌檔案嗎？此操作不可恢復。'
    }
  },
  
  mappings: {
    regex: '正則表達式',
    replacement: '替換內容',
    regexPlaceholder: '^/雲端硬碟/(.*)$',
    replacementPlaceholder: '/mnt/gd/$1',
    empty: '暫無對映規則',
    addFirst: '新增第一條規則'
  },
  
  cron: {
    placeholder: '0 0 * * *',
    invalid: '無效的 Cron 表達式',
    everyHour: '每小時',
    everyMinute: '每分鐘',
    atTime: '{hour}:{minute}',
    custom: '自訂時間',
    daily: '每天',
    weekly: '每週',
    monthly: '每月',
    fields: {
      minute: '分',
      hour: '時',
      day: '日',
      month: '月',
      weekday: '週'
    }
  },
  
  panels: {
    basic: {
      title: '基礎設定',
      description: '設定伺服器連接埠、認證和 SSL 設定',
      server: '伺服器設定',
      port: '監聽連接埠',
      portHint: '服務監聽連接埠，預設 8448',
      publicUrl: '公網地址',
      publicUrlPlaceholder: 'https://your-domain.com',
      publicUrlHint: '用於 Google Drive Webhook 回呼',
      webhookPath: 'Webhook 路徑',
      auth: '存取認證',
      username: '使用者名稱',
      password: '密碼',
      ssl: 'SSL 設定',
      enableSsl: '啟用 HTTPS',
      sslCert: '憑證路徑',
      sslCertPlaceholder: '/path/to/cert.pem',
      sslKey: '私鑰路徑',
      sslKeyPlaceholder: '/path/to/key.pem'
    },
    
    integrations: {
      title: '服務整合',
      description: '設定 Google Drive、Rclone 和 Symedia 整合',
      googleDrive: 'Google Drive',
      qps: 'API QPS',
      qpsHint: '每秒 API 請求限制',
      personalDriveName: '個人雲端硬碟名稱',
      personalDriveNamePlaceholder: '我的雲端硬碟',
      rclone: 'Rclone 實例',
      rcloneHost: '主機地址',
      rcloneEndpoint: 'API 端點',
      waitForData: '等待資料同步',
      noRcloneInstances: '未設定 Rclone 實例',
      addRcloneInstance: '新增實例',
      symedia: 'Symedia 通知',
      symediaHost: '主機地址',
      symediaEndpoint: 'API 端點',
      symediaTemplate: '請求範本',
      symediaTemplateHint: '可用變數：',
      notifyUnmatched: '通知未匹配的路徑',
      headers: '請求頭',
      noHeaders: '未設定請求頭'
    },
    
    mappings: {
      title: '路徑對映',
      description: '設定 SA 和 Rclone 的路徑轉換規則',
      saMappings: 'SA 路徑對映',
      saMappingsDesc: '將 Google Drive 路徑轉換為 Symedia 可識別的路徑',
      rcloneMappings: 'Rclone 路徑對映',
      rcloneMappingsDesc: '將 Google Drive 路徑轉換為 Rclone VFS 路徑'
    },
    
    ignore: {
      title: '忽略設定',
      description: '設定需要忽略的 Google Drive 資料夾',
      placeholder: '輸入 Google Drive 資料夾 ID',
      empty: '未設定忽略規則',
      emptyHint: '新增資料夾 ID 以忽略其中的檔案變更'
    },
    
    advanced: {
      title: '進階設定',
      description: '設定同步參數和日誌管理',
      syncSettings: '同步設定',
      debounce: '防抖延遲',
      debounceHint: '收到變更通知後等待的時間，避免頻繁觸發',
      logging: '日誌設定',
      logDir: '日誌目錄',
      logLevel: '日誌層級',
      logLevels: {
        quiet: '靜默',
        info: '資訊',
        debug: '偵錯'
      },
      enableLogSave: '儲存日誌到檔案',
      logCleanup: '日誌清理',
      enableLogCleanup: '啟用自動清理',
      retentionDays: '保留天數',
      cleanupCron: '清理排程'
    },
    
    oauth: {
      title: 'OAuth 認證',
      description: '設定 Google OAuth 憑證',
      credentials: 'OAuth 憑證',
      clientId: 'Client ID',
      clientIdPlaceholder: '輸入 Google OAuth Client ID',
      clientSecret: 'Client Secret',
      clientSecretPlaceholder: '輸入 Client Secret',
      redirectUri: '重新導向 URI',
      redirectUriPlaceholder: 'http://localhost:8448/api/oauth/callback',
      redirectUriHint: '需要在 Google Cloud Console 中設定',
      authorization: '授權操作',
      actionDesc: '設定好 OAuth 憑證後，點擊下方按鈕跳轉到 Google 授權頁面',
      goToGoogle: '前往 Google 授權',
      urlOpened: '已在新視窗開啟授權頁面',
      urlError: '取得授權連結失敗，請檢查 OAuth 設定',
      noClientId: '請先填寫 Client ID'
    }
  },
  
  settings: {
    theme: '主題',
    themeLight: '淺色',
    themeDark: '深色',
    themeSystem: '跟隨系統',
    language: '語言'
  },
  
  footer: {
    version: '版本',
    docs: '文件',
    github: 'GitHub'
  }
}
