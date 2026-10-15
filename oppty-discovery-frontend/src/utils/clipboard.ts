/**
 * 健壮的跨环境剪贴板复制工具函数
 * 兼容在非 HTTPS（非安全上下文，如内网 IP 访问 http://192.168.8.206）环境下的文本复制
 *
 * @param text 需要复制的纯文本或链接字符串
 * @returns Promise<boolean> 是否复制成功
 */
export const copyToClipboard = async (text: string): Promise<boolean> => {
  if (!text) {
    return false;
  }

  // 1. 若支持现代 Clipboard API 且处于安全上下文 (HTTPS 或 localhost) 中，优先使用异步 API
  if (typeof navigator !== 'undefined' && navigator.clipboard && window.isSecureContext) {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch (e) {
      console.warn('[剪贴板] 现代 API 写入受限，自动降级至传统复制方式:', e);
    }
  }

  // 2. 传统降级模式：动态创建隐藏 textarea 元素并触发 document.execCommand('copy')
  try {
    const textArea = document.createElement('textarea');
    textArea.value = text;
    // 使用 fixed 定位和屏幕外绝对坐标，避免在移动端或大屏幕引起页面视口滚动跳动
    textArea.style.position = 'fixed';
    textArea.style.top = '-9999px';
    textArea.style.left = '-9999px';
    textArea.style.opacity = '0';
    textArea.setAttribute('readonly', '');
    document.body.appendChild(textArea);
    textArea.select();
    textArea.setSelectionRange(0, text.length);

    const successful = document.execCommand('copy');
    document.body.removeChild(textArea);
    if (!successful) {
      console.warn('[剪贴板警告] document.execCommand 降级复制返回失败');
    }
    return successful;
  } catch (err) {
    console.error('[剪贴板错误] 降级复制执行失败:', err);
    return false;
  }
};
