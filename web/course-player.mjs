export function createCoursePlayer(root, repository) {
  const find = (selector) => root.querySelector(selector);
  const make = (tag, className, text) => {
    const node = document.createElement(tag); if(className)node.className=className;
    if(text!==undefined)node.textContent=text; return node;
  };
  let lab, frames=[], index=-1, timer=null;
  function stop() {
    if(timer!==null)clearInterval(timer);timer=null;
    find('[data-action="play"]').textContent=index===frames.length-1?'重新播放':'播放过程';
  }
  function render() {
    const current=frames[index];
    find('.course-diagram').replaceChildren();
    lab.nodes.forEach((name,i)=>{
      const node=make('div',`concept-node${current?.node===i?' current':''}${current?.node===i&&current.failed?' failed':''}`);
      node.append(make('span','concept-number',String(i+1).padStart(2,'0')),make('strong',null,name));
      if(current?.node===i)node.append(make('span','concept-marker',current.failed?'此处返回 / 拒绝':'当前步骤'));
      find('.course-diagram').append(node);
    });
    find('.course-frame-message').textContent=current?.message || '选择场景，点击播放或下一步，观察本课代码的执行过程。';
    find('.course-frame-code').textContent=current?.code || '// 选择场景后开始';
    find('.course-values').replaceChildren();
    if(!current)find('.course-values').append(make('p','value-empty','变量与输出会随步骤更新。'));
    else for(const [key,value] of Object.entries(current.values)){
      const pair=make('div');pair.append(make('dt',null,key),make('dd',null,String(value)));find('.course-values').append(pair);
    }
    find('.course-trace').replaceChildren();
    frames.slice(0,index+1).forEach((f)=>find('.course-trace').append(make('li',null,f.message)));
    find('.course-frame-count').textContent=`${index+1} / ${frames.length}`;
    find('[data-action="next"]').disabled=index===frames.length-1;
    find('[data-action="previous"]').disabled=index<0;
  }
  function advance() {if(index<frames.length-1)index++;render();if(index===frames.length-1)stop();}
  function reset() {stop();index=-1;frames=lab.steps(find('.course-scenario').value);render();find('[data-action="play"]').textContent='播放过程';}
  find('.course-scenario').addEventListener('change',reset);
  find('[data-action="reset"]').addEventListener('click',reset);
  find('[data-action="next"]').addEventListener('click',()=>{stop();advance();});
  find('[data-action="previous"]').addEventListener('click',()=>{stop();index--;render();});
  find('[data-action="play"]').addEventListener('click',()=>{
    if(timer!==null){stop();return;}
    if(index===frames.length-1)index=-1;
    advance();if(index===frames.length-1)return;
    find('[data-action="play"]').textContent='暂停';timer=setInterval(advance,1200);
  });
  return {stop,show(next){
    stop();lab=next;
    find('.course-scenario').replaceChildren();
    lab.options.forEach(({id,label})=>{const option=make('option',null,label);option.value=id;find('.course-scenario').append(option);});
    find('.course-lab-note').textContent=lab.note;
    find('.course-lab-source').href=`${repository}/blob/main/${lab.source}`;
    find('.course-lab-source').textContent=`对照 ${lab.source} ↗`;
    root.querySelector('details').open=false;
    reset();
  }};
}
