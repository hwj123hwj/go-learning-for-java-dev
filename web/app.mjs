import { labForCourse } from './course-labs.mjs';
import { createCoursePlayer } from './course-player.mjs';
import { channelState, channelAction, gatewaySteps } from './models.mjs';
import { courses, repository } from './courses.mjs';
const $ = (selector) => document.querySelector(selector);
const make = (tag, className, text) => {
  const node = document.createElement(tag); if (className) node.className = className;
  if (text !== undefined) node.textContent = text; return node;
};
let completed = new Set();
try {
  const saved = JSON.parse(localStorage.getItem('java-go-progress-v1') || '[]');
  if (Array.isArray(saved)) completed = new Set(saved.filter((id) => courses.some((course) => course.id === id)));
} catch { $('#storage-note').textContent = '浏览器存储不可用或数据异常，当前进度仍可在本页使用。'; }
let active = '01';
const coursePlayer = createCoursePlayer($('#course-panel'), repository);
let initialized = false;
function renderProgress() {
  $('#progress-count').textContent = `${completed.size} / ${courses.length}`;
  $('#course-progress').value = completed.size;
  courses.forEach((course) => {
    const item = $(`[data-course="${course.id}"]`);
    item.setAttribute('aria-current', course.id === active ? 'step' : 'false');
    item.querySelector('.course-done').textContent = completed.has(course.id) ? '✓' : '';
    item.setAttribute('aria-label', `阶段 ${course.id} ${course.name}${completed.has(course.id) ? '，已完成' : ''}`);
  });
  $('#lesson-complete').checked = completed.has(active);
}
function saveProgress() {
  try { localStorage.setItem('java-go-progress-v1', JSON.stringify([...completed])); }
  catch { $('#storage-note').textContent = '进度未能保存，刷新后可能丢失。'; }
  renderProgress();
}
courses.forEach((course) => {
  const item = make('a', 'course-item'); item.href = `#course-${course.id}`; item.dataset.course = course.id;
  const words = make('span'); words.append(make('strong', null, course.name), make('small', null, course.sub));
  const done = make('span', 'course-done'); done.setAttribute('aria-hidden', 'true');
  item.append(make('span', 'course-number', course.id), words, done); $('#course-list').append(item);
});
function selectCourse() {
  const id = location.hash.match(/^#course-(\d{2})$/)?.[1];
  if (!id && location.hash && initialized) return;
  const course = courses.find((c) => c.id === id) || courses[0]; active = course.id;
  initialized = true;
  $('#lesson-stage').textContent = `阶段 ${course.id} · ${course.name}`;
  $('#lesson-title').textContent = course.title; $('#lesson-description').textContent = course.description;
  $('#java-bridge').textContent = course.bridge; $('#lesson-goal').textContent = course.goal;
  $('#lesson-command').textContent = course.command; $('#copy-status').textContent = '';
  $('#lesson-source').href = `${repository}/tree/main/${course.directory}`;
  $('#lesson-exercises').href = `${repository}/${course.exercise ? 'blob' : 'tree'}/main/${course.exercise || course.directory}`;
  $('#lesson-exercises').textContent = course.exercise?.includes('EXERCISES') ? '进入练习 ↗' : course.exercise ? '阅读学习说明 ↗' : '查看服务代码 ↗';
  $('#solution-link').hidden = !course.solution;
  if (course.solution) $('#solution-link').href = `${repository}/blob/main/${course.solution}`;
  showCourseLab(active);
  document.title = `${course.name} · Java → Go 交互学习实验室`;
  renderProgress();
  revealActiveCourse();
}
function revealActiveCourse() {
  if (!window.matchMedia('(max-width: 760px)').matches) return;
  const list = $('#course-list'); const item = $(`[data-course="${active}"]`);
  list.scrollLeft = item.offsetLeft - list.offsetLeft - (list.clientWidth - item.clientWidth) / 2;
}
window.matchMedia('(max-width: 760px)').addEventListener('change', revealActiveCourse);
window.addEventListener('pageshow', () => requestAnimationFrame(revealActiveCourse));
$('#lesson-complete').addEventListener('change', (event) => {
  if (event.target.checked) completed.add(active); else completed.delete(active); saveProgress();
});
$('#clear-progress').addEventListener('click', () => { completed.clear(); saveProgress(); });
$('#copy-command').addEventListener('click', async () => {
  try { await navigator.clipboard.writeText($('#lesson-command').textContent); $('#copy-status').textContent = '已复制。请从仓库根目录执行。'; }
  catch { const range = document.createRange(); range.selectNodeContents($('#lesson-command')); const selection = window.getSelection(); selection.removeAllRanges(); selection.addRange(range); $('#copy-status').textContent = '自动复制不可用，命令已选中，请手动复制。'; }
});
let channel = channelState(2);
function renderChannel() {
  $('#buffer').replaceChildren();
  if (!channel.capacity) $('#buffer').append(make('span', 'handoff', '无缓冲 · 同步交接'));
  for (let i = 0; i < channel.capacity; i++) {
    const slot = make('span', `slot${i < channel.queue.length ? ' filled' : ''}`, channel.queue[i] ?? '·');
    slot.setAttribute('aria-label', i < channel.queue.length ? `槽位 ${i + 1}：${channel.queue[i]}` : `槽位 ${i + 1}：空`);
    $('#buffer').append(slot);
  }
  $('.buffer-label code').textContent = channel.capacity ? `ch := make(chan int, ${channel.capacity})` : 'ch := make(chan int)';
  $('#buffer-count').textContent = `${channel.queue.length} / ${channel.capacity} 已使用${channel.closed ? ' · 已关闭' : ''}`;
  $('#channel-kind').textContent = channel.capacity ? `缓冲区可容纳 ${channel.capacity} 个值` : '发送和接收必须配对';
  $('#sender-value').textContent = channel.sender !== null ? `等待发送 ${channel.sender}` : channel.closed ? '已停止发送' : `准备发送 ${channel.next}`;
  $('#sender-state').textContent = channel.sender !== null ? '阻塞 · 等待接收' : channel.closed ? 'Channel 已关闭' : '就绪';
  $('#sender-state').className = channel.sender !== null ? 'state waiting' : 'state';
  $('#receiver-value').textContent = channel.received ? `值 ${channel.received.value} · ok=${channel.received.ok}` : '尚未接收';
  $('#receiver-state').textContent = channel.receiverWaiting ? '阻塞 · 等待发送' : '就绪';
  $('#receiver-state').className = channel.receiverWaiting ? 'state waiting' : 'state';
  $('#send').disabled = channel.sender !== null || channel.closed;
  $('#receive').disabled = channel.receiverWaiting;
  $('#channel-close').disabled = channel.closed || channel.sender !== null;
  $('#channel-message').textContent = channel.message;
  $('#close-note').textContent = channel.sender !== null ? '先接收以解除发送阻塞。关闭有待发送者的 Channel 会导致发送方 panic，本实验禁止这一操作。' : '由发送方关闭；关闭后可以继续读完缓冲值。';
}
for (const [selector, action] of [['#send','send'], ['#receive','receive'], ['#channel-close','close']]) {
  $(selector).addEventListener('click', () => { channel = channelAction(channel, action); renderChannel(); });
}
function resetChannel() { channel = channelState(Number($('#capacity').value)); renderChannel(); }
$('#capacity').addEventListener('change', resetChannel); $('#channel-reset').addEventListener('click', resetChannel);
let frames = gatewaySteps('success'); let frame = -1; let timer = null;
function pause() { if (timer !== null) clearInterval(timer); timer = null; $('#play').textContent = frame === frames.length - 1 ? '重新播放' : '播放过程'; }
function renderGateway() {
  const current = frames[frame];
  document.querySelectorAll('[data-node]').forEach((node) => {
    node.classList.toggle('active', node.dataset.node === current?.node);
    node.classList.toggle('failed', node.dataset.node === current?.failed);
  });
  $('#gateway-message').textContent = current?.message || '选择场景后播放，或逐步查看。';
  $('#stream-text').textContent = current?.output || (frame < 0 ? '等待请求…' : '尚未收到内容');
  $('#stream-state').textContent = current?.state || '尚未发送';
  $('#step-count').textContent = `${frame + 1} / ${frames.length}`;
  $('#step').disabled = frame === frames.length - 1;
  $('#trace').replaceChildren();
  if (frame < 0) $('#trace').append(make('li', null, '点击“播放过程”开始。'));
  else frames.slice(0, frame + 1).forEach((s) => $('#trace').append(make('li', null, s.event)));
}
function advance() { if (frame < frames.length - 1) frame++; renderGateway(); if (frame === frames.length - 1) pause(); }
function resetGateway() { pause(); frame = -1; frames = gatewaySteps($('#scenario').value); renderGateway(); $('#play').textContent = '播放过程'; }
$('#scenario').addEventListener('change', resetGateway); $('#gateway-reset').addEventListener('click', resetGateway);
$('#step').addEventListener('click', () => { pause(); advance(); });
$('#play').addEventListener('click', () => {
  if (timer !== null) { pause(); return; }
  if (frame === frames.length - 1) frame = -1;
  advance(); if (frame === frames.length - 1) return;
  $('#play').textContent = '暂停'; timer = setInterval(advance, 1200);
});
function showCourseLab(id) {
  pause(); coursePlayer.stop();
  const selection = labForCourse(id);
  $('#lab-title').textContent = selection.title;
  for (const name of ['channel', 'gateway', 'course']) $(`#${name}-panel`).hidden = name !== selection.panel;
  if (selection.panel === 'course') coursePlayer.show(selection.lab);
  else if (selection.panel === 'channel') resetChannel();
  else resetGateway();
}
document.addEventListener('visibilitychange', () => { if (document.hidden) { pause(); coursePlayer.stop(); } });
window.addEventListener('pagehide', () => { pause(); coursePlayer.stop(); }); window.addEventListener('hashchange', selectCourse);
document.addEventListener('keydown', (event) => {
  if (event.key !== 'ArrowLeft' && event.key !== 'ArrowRight') return;
  if (event.metaKey || event.ctrlKey || event.altKey || event.shiftKey) return;
  if (event.target instanceof Element && event.target.closest('input, select, textarea, [contenteditable]')) return;
  const panel = document.querySelector('#course-panel:not([hidden]), #gateway-panel:not([hidden])');
  if (!panel) return;
  const next = event.key === 'ArrowRight';
  const button = panel.id === 'course-panel' ? panel.querySelector(`[data-action="${next ? 'next' : 'previous'}"]`) : next ? panel.querySelector('#step') : null;
  if (button && !button.disabled) { event.preventDefault(); button.click(); }
});
renderChannel(); renderGateway(); selectCourse();
