async function loadStats(){
  document.getElementById('stats').textContent = JSON.stringify(await (await fetch('/api/stats')).json(), null, 2);
  document.getElementById('series').textContent = JSON.stringify(await (await fetch('/api/series')).json(), null, 2);
}
document.getElementById('append').onclick = async () => {
  const body = {
    series: document.getElementById('name').value,
    labels: {host: document.getElementById('host').value},
    value: Number(document.getElementById('value').value)
  };
  const r = await fetch('/api/append', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(body)});
  document.getElementById('msg').textContent = r.ok ? 'ok' : await r.text();
  loadStats();
};
document.getElementById('compact').onclick = async () => {
  const r = await fetch('/api/compact', {method:'POST'});
  document.getElementById('msg').textContent = r.ok ? 'compacted' : await r.text();
  loadStats();
};
document.getElementById('query').onclick = async () => {
  const s = document.getElementById('qseries').value;
  const data = await (await fetch('/api/query?series=' + encodeURIComponent(s))).json();
  document.getElementById('result').textContent = JSON.stringify(data, null, 2);
};
loadStats();
setInterval(loadStats, 3000);
