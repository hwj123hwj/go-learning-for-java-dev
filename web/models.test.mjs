import test from 'node:test';
import assert from 'node:assert/strict';
import { existsSync } from 'node:fs';
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
