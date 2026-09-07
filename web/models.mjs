// Deterministic teaching models; these do not execute Go or make network requests.
export function channelState(capacity = 2) {
  return { capacity, queue: [], next: 1, sender: null, receiverWaiting: false, received: null, closed: false,
    message: capacity ? `先发送 ${capacity} 次，再试着多发一次。缓冲区满时会发生什么？` : '没有缓冲区。试着先发送，再接收，观察一次同步交接。' };
}
export function channelAction(previous, action) {
  const s = { ...previous, queue: [...previous.queue] };
  if (action === 'send') {
    if (s.closed || s.sender !== null) return s;
    const value = s.next++;
    if (s.receiverWaiting) {
      s.receiverWaiting = false; s.received = { value, ok: true };
      s.message = `值 ${value} 直接交给等待中的接收者，两侧继续执行。`;
    } else if (s.queue.length < s.capacity) {
      s.queue.push(value); s.message = `值 ${value} 进入缓冲区。发送完成，不必等待接收者。`;
    } else {
      s.sender = value;
      s.message = s.capacity ? `缓冲区已满，发送者阻塞在值 ${value} 上。接收一次才能继续。` : `发送者阻塞在值 ${value} 上，等待接收者配对。`;
    }
  }
  if (action === 'receive') {
    if (s.receiverWaiting) return s;
    if (s.queue.length) {
      const value = s.queue.shift(); s.received = { value, ok: true };
      if (s.sender !== null) {
        const pending = s.sender; s.queue.push(pending); s.sender = null;
        s.message = `接收到 ${value}；空出的槽位接纳等待发送的 ${pending}，发送者恢复。`;
      } else s.message = `接收到 ${value}，ok=true。${s.closed ? '关闭后，剩余缓冲值仍然可读。' : '缓冲区腾出一个位置。'}`;
    } else if (s.sender !== null) {
      const value = s.sender; s.received = { value, ok: true }; s.sender = null;
      s.message = `值 ${value} 完成交接，两侧恢复。无缓冲 Channel 不储存值。`;
    } else if (s.closed) {
      s.received = { value: 0, ok: false }; s.message = 'Channel 已关闭且已读空：立即返回 int 零值 0，ok=false，不再阻塞。';
    } else {
      s.receiverWaiting = true; s.message = '当前没有可接收的值，接收者阻塞，等待发送者。';
    }
  }
  if (action === 'close' && !s.closed && s.sender === null) {
    s.closed = true;
    if (s.receiverWaiting) { s.receiverWaiting = false; s.received = { value: 0, ok: false }; }
    s.message = '发送方关闭 Channel。不能再发送；缓冲值读完后，接收得到 0 和 ok=false。';
  }
  return s;
}
const start = [
  { node: 'client', message: '客户端提交 stream=true 的请求。', event: 'client → proxy', output: '', state: '请求已发送' },
  { node: 'proxy', message: '代理创建请求 Context，依次尝试配置的上游。', event: 'Chain.Stream → A', output: '', state: '等待上游' },
];
const stream = (node, name) => [
  { node, message: `${name} 返回成功响应，准备读取 SSE。`, event: `${name}: HTTP 200`, output: '', state: '等待首个片段' },
  { node: 'proxy', message: '解析首个 delta，立即 Flush 给客户端。', event: 'delta → Flush', output: '你好，', state: '流式输出中' },
  { node: 'client', message: '继续逐块转发，无需等待完整回答。', event: 'delta → Flush', output: '你好，Go 开发者。', state: '流式输出中' },
  { node: 'client', message: '上游 [DONE]，代理写出结束片段并完成请求。', event: '[DONE] → 完成', output: '你好，Go 开发者。', state: '已完成' },
];
export function gatewaySteps(scenario) {
  if (scenario === 'fallback') return [...start,
    { node: 'a', failed: 'a', message: 'A 返回 503，可切换错误；流式路径不在 A 重试，直接尝试 B。', event: 'A: 503 → fallback B', output: '', state: '正在切换上游' },
    ...stream('b', 'B').map((s) => ({ ...s, failed: 'a' }))];
  if (scenario === 'invalid') return [...start,
    { node: 'a', failed: 'a', message: 'A 返回 400，属于不可重试错误。停止链路，不调用 B。', event: 'A: 400 → error', output: '', state: '请求失败 · 不切换' }];
  if (scenario === 'partial') return [...start, ...stream('a', 'A').slice(0, 2),
    { node: 'a', failed: 'a', message: '读流返回错误。客户端已经看到部分内容，不能重放或切换 B；返回流内错误。', event: 'read error → stream error', output: '你好，', state: '回答不完整 · 不切换' }];
  if (scenario === 'cancel') return [...start, ...stream('a', 'A').slice(0, 2),
    { node: 'client', message: '客户端断开，请求 Context 取消，上游读取结束。已显示的内容保留，停止输出。', event: 'ctx canceled → stop', output: '你好，', state: '已取消' }];
  return [...start, ...stream('a', 'A')];
}
