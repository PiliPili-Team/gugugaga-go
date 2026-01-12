/**
 * Dialog Utility - Custom alert and confirm dialogs
 */

export function showAlert(title: string, message: string): Promise<void> {
  return new Promise((resolve) => {
    const event = new CustomEvent('show-alert', {
      detail: { title, message, resolve }
    })
    window.dispatchEvent(event)
    
    // Wait for dialog to close
    const checkClosed = setInterval(() => {
      // This will be handled by Dashboard component
      resolve()
      clearInterval(checkClosed)
    }, 100)
  })
}

export function showConfirm(title: string, message: string): Promise<boolean> {
  return new Promise((resolve) => {
    const event = new CustomEvent('show-confirm', {
      detail: { title, message }
    })
    window.dispatchEvent(event)
    
    // Listen for result
    const handleResult = ((e: CustomEvent) => {
      resolve(e.detail)
      window.removeEventListener('confirm-result', handleResult as EventListener)
    }) as EventListener
    
    window.addEventListener('confirm-result', handleResult)
  })
}
