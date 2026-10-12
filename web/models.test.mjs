import test from 'node:test';
import assert from 'node:assert/strict';
import { existsSync, readFileSync } from 'node:fs';
import { channelState, channelAction, gatewaySteps } from './models.mjs';
import { courses } from './courses.mjs';
const act = (s, ...actions) => actions.reduce(channelAction, s);
test('unbuffered send waits for receiver and cannot issue a second send', () => {
  const initial = channelState(0), waiting = act(initial, 'send', 'send');
  assert.equal(waiting.sender, 1); assert.equal(waiting.next, 2); assert.deepEqual(initial.queue, []);
  const done = act(waiting, 'receive'); assert.deepEqual(done.received, {value:1,ok:true}); assert.equal(done.sender, null);
});
test('waiting receiver rendezvous bypasses buffer', () => {
  const state = act(channelState(2), 'receive', 'send');
  assert.equal(state.receiverWaiting,false); assert.deepEqual(state.queue,[]); assert.equal(state.received.value,1);
});
test('buffered FIFO replenishes from a blocked sender', () => {
  const state = act(channelState(2), 'send','send','send','receive');
  assert.deepEqual(state.queue,[2,3]); assert.equal(state.received.value,1); assert.equal(state.sender,null);
});
test('close drains existing values before returning zero and rejects sends', () => {
  const closed = act(channelState(2),'send','send','close','send');
  assert.equal(closed.next,3); assert.equal(closed.closed,true);
  const one = act(closed,'receive'), two = act(one,'receive'), zero = act(two,'receive');
  assert.deepEqual(one.received,{value:1,ok:true}); assert.deepEqual(two.received,{value:2,ok:true});
  assert.deepEqual(zero.received,{value:0,ok:false});
});
test('closing wakes a waiting receiver; pending send prevents closing in this demo', () => {
  assert.deepEqual(act(channelState(0),'receive','close').received,{value:0,ok:false});
  assert.equal(act(channelState(0),'send','close').closed,false);
});
test('503 switches to B before output; 400 never switches', () => {
  const steps = gatewaySteps('fallback'); const fallback = steps.findIndex(s=>s.event.includes('fallback'));
  assert.ok(fallback>=0); assert.ok(steps.slice(0,fallback+1).every(s=>s.output===''));
  assert.ok(steps.some(s=>s.node==='b')); assert.equal(steps.at(-1).state,'已完成');
  assert.ok(gatewaySteps('invalid').every(s=>s.node!=='b'));
});
test('partial failure and cancel preserve output and never call B', () => {
  for (const scenario of ['partial','cancel']) {
    const steps=gatewaySteps(scenario); assert.equal(steps.at(-1).output,'你好，');
    assert.ok(steps.every(s=>s.node!=='b')); assert.ok(steps.every(s=>s.event!=='[DONE] → 完成'));
  }
});
test('course content links resolve to actual repository files', () => {
  assert.equal(courses.length,9); assert.equal(new Set(courses.map(c=>c.id)).size,9);
  for (const course of courses) for(const path of [course.directory, course.exercise, course.solution].filter(Boolean)) {
    assert.ok(existsSync(new URL(`../${path}`,import.meta.url)),`missing course target ${path}`);
  }
});
test('share image and PWA icons exist and are wired up', () => {
  const manifest = JSON.parse(readFileSync(new URL('./manifest.webmanifest', import.meta.url), 'utf8'));
  const html = readFileSync(new URL('./index.html', import.meta.url), 'utf8');
  for (const icon of manifest.icons) assert.ok(existsSync(new URL(`./${icon.src}`, import.meta.url)), `missing icon ${icon.src}`);
  assert.ok(existsSync(new URL('./og-card.png', import.meta.url)), 'missing og-card.png');
  assert.ok(html.includes('og-card.png')); assert.ok(html.includes('icon-192.png'));
});

import { courseLabs, labForCourse } from './course-labs.mjs';
test('every course has its own lab and every scenario has valid frames and source', () => {
  assert.equal(new Set(courses.map(c => labForCourse(c.id).title)).size, 9);
  assert.equal(labForCourse('06').panel, 'channel');
  assert.equal(labForCourse('09').panel, 'gateway');
  for (const [id, lab] of Object.entries(courseLabs)) {
    assert.equal(labForCourse(id).panel, 'course');
    assert.ok(existsSync(new URL(`../${lab.source}`, import.meta.url)));
    for (const option of lab.options) {
      const frames = lab.steps(option.id);
      assert.ok(frames.length > 1);
      for (const frame of frames) {
        assert.ok(frame.node >= 0 && frame.node < lab.nodes.length);
        assert.ok(frame.message && frame.code && Object.keys(frame.values).length);
      }
    }
  }
});
const last = (id, option) => courseLabs[id].steps(option).at(-1).values;
test('division errors and integer results remain distinct', () => {
  assert.equal(last('01','zero').result, 0);
  assert.notEqual(last('01','zero').err, 'nil');
  assert.equal(last('01','integer').result, 3);
});
test('invalid API IDs and failed authentication stop before business work', () => {
  assert.equal(last('02','invalid').HTTP, 400);
  assert.equal(last('02','invalid')['数据库查询次数'], 0);
  assert.ok(courseLabs['02'].steps('invalid').every(f => f.node !== 2 && f.node !== 3));
  assert.equal(last('02','missing').HTTP, 404);
  assert.equal(last('02','found').HTTP, 200);
  for(const option of ['missing','malformed','expired']) {
    assert.equal(last('03',option).HTTP, 401);
    assert.ok(courseLabs['03'].steps(option).every(f => f.node !== 3));
  }
  assert.equal(last('03','valid')['业务处理'], '放行');
});
test('RPC message counts distinguish all four modes', () => {
  for(const [option, sent, received] of [['unary',1,1],['server',1,3],['client',2,1],['bidi',3,3]]) {
    assert.equal(last('04',option)['客户端消息数'], sent);
    assert.equal(last('04',option)['服务端消息数'], received);
  }
});
test('downstream creation errors never write orders; reads can degrade', () => {
  for(const option of ['missing','cancel']) assert.equal(last('05',option)['订单写入次数'], 0);
  assert.equal(last('05','degrade').user_name, 'unknown');
});
test('token bucket enforces burst capacity and replenishes with virtual time', () => {
  for(const [option, pass, reject] of [['burst',3,3],['paced',6,0],['fast',4,2]]) {
    assert.equal(last('07',option)['通过'], pass);
    assert.equal(last('07',option)['拒绝'], reject);
    for(const frame of courseLabs['07'].steps(option)) assert.ok(frame.values['令牌'] >= 0 && frame.values['令牌'] <= 3);
  }
});
test('error wrapping preserves identity while string formatting loses it', () => {
  assert.equal(last('08','wrap')['匹配'], true);
  assert.equal(last('08','string')['匹配'], false);
  assert.equal(last('08','typed')['ve.Field'], 'age');
});
