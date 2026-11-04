using System;
using System.ComponentModel;
using Avalonia.Threading;

namespace XenAdmin.Desktop
{
    /// <summary>Lets XenModel's InvokeHelper marshal cache updates onto the Avalonia UI thread.</summary>
    public class DispatcherSynchronizer : ISynchronizeInvoke
    {
        public bool InvokeRequired => !Dispatcher.UIThread.CheckAccess();

        public IAsyncResult BeginInvoke(Delegate method, object[] args)
        {
            return Dispatcher.UIThread.InvokeAsync(() => method.DynamicInvoke(args)).GetTask();
        }

        public object EndInvoke(IAsyncResult result)
        {
            return ((System.Threading.Tasks.Task<object>)result).GetAwaiter().GetResult();
        }

        public object Invoke(Delegate method, object[] args)
        {
            return Dispatcher.UIThread.Invoke(() => method.DynamicInvoke(args));
        }
    }
}
