import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import ForceGraph2D from 'react-force-graph-2d';

const NODE_TYPES = {
  NETWORK: 0,
  HOST: 1,
  GATEWAY: 2,
  UNKNOWN: 3
};

const NODE_COLORS = {
  [NODE_TYPES.NETWORK]: '#8fd3ff',
  [NODE_TYPES.HOST]: '#d2b6ff',
  [NODE_TYPES.GATEWAY]: '#ffd166',
  [NODE_TYPES.UNKNOWN]: '#9aa4b2'
};

const NODE_LABELS = {
  [NODE_TYPES.NETWORK]: 'Network',
  [NODE_TYPES.HOST]: 'Host',
  [NODE_TYPES.GATEWAY]: 'Gateway',
  [NODE_TYPES.UNKNOWN]: 'Unknown'
};

const DEFAULT_GRAPH = { nodes: [], links: [] };

const linkColorFor = (link) => {
  switch (link.type) {
    case 'routes-via':
      return 'rgba(255, 209, 102, 0.7)';
    case 'member-of':
      return 'rgba(143, 211, 255, 0.5)';
    case 'responds-to':
      return 'rgba(210, 182, 255, 0.6)';
    default:
      return 'rgba(154, 164, 178, 0.5)';
  }
};

const formatProtocols = (protocols) => {
  if (!protocols) return 'None';
  return Object.keys(protocols)
    .filter((key) => protocols[key])
    .map((key) => key.toUpperCase())
    .join(', ');
};

function App() {
  const [graphData, setGraphData] = useState(DEFAULT_GRAPH);
  const [status, setStatus] = useState('idle');
  const [lastUpdated, setLastUpdated] = useState(null);
  const graphRef = useRef(null);

  const fetchGraph = useCallback(async () => {
    setStatus('loading');
    try {
      const response = await fetch('/api/graph');
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }
      const data = await response.json();
      setGraphData(data ?? DEFAULT_GRAPH);
      setLastUpdated(new Date());
      setStatus('success');
    } catch (error) {
      setStatus('error');
    }
  }, []);

  useEffect(() => {
    fetchGraph();
    const interval = setInterval(fetchGraph, 10000);
    return () => clearInterval(interval);
  }, [fetchGraph]);

  useEffect(() => {
    if (graphRef.current && graphData.nodes.length > 0) {
      graphRef.current.zoomToFit(400, 80);
    }
  }, [graphData]);

  const legendItems = useMemo(
    () =>
      Object.entries(NODE_LABELS).map(([key, label]) => ({
        key,
        label,
        color: NODE_COLORS[key]
      })),
    []
  );

  return (
    <div className="app">
      <header className="toolbar">
        <div>
          <h1>NetMap Graph</h1>
          <p className="subtitle">Live topology view powered by the NetMap HTTP service.</p>
        </div>
        <div className="actions">
          <button type="button" onClick={fetchGraph} disabled={status === 'loading'}>
            {status === 'loading' ? 'Refreshing…' : 'Refresh'}
          </button>
          <div className="status">
            <span className={`status-dot status-${status}`} />
            <span>
              {status === 'error'
                ? 'Service unreachable'
                : status === 'loading'
                ? 'Fetching graph…'
                : 'Connected'}
            </span>
          </div>
        </div>
      </header>

      <div className="content">
        <aside className="sidebar">
          <section>
            <h2>Legend</h2>
            <ul>
              {legendItems.map((item) => (
                <li key={item.key}>
                  <span className="legend-dot" style={{ background: item.color }} />
                  {item.label}
                </li>
              ))}
            </ul>
          </section>
          <section>
            <h2>Graph Summary</h2>
            <p>
              Nodes: <strong>{graphData.nodes.length}</strong>
            </p>
            <p>
              Links: <strong>{graphData.links.length}</strong>
            </p>
            <p>
              Last updated:{' '}
              <strong>
                {lastUpdated ? lastUpdated.toLocaleTimeString() : 'Not yet loaded'}
              </strong>
            </p>
          </section>
          <section>
            <h2>Tips</h2>
            <ul>
              <li>Drag nodes to pin positions.</li>
              <li>Scroll to zoom, click and drag background to pan.</li>
              <li>Hover nodes for IP, vendor, and protocol details.</li>
            </ul>
          </section>
        </aside>

        <main className="graph-panel">
          <ForceGraph2D
            ref={graphRef}
            graphData={graphData}
            backgroundColor="#0b0f1c"
            nodeRelSize={6}
            linkDirectionalParticles={2}
            linkDirectionalParticleWidth={2}
            linkDirectionalParticleSpeed={0.004}
            linkColor={(link) => linkColorFor(link)}
            linkWidth={(link) => (link.type === 'routes-via' ? 2 : 1)}
            nodeCanvasObject={(node, ctx, globalScale) => {
              const label = node.ip || node.id;
              const fontSize = 12 / globalScale;
              ctx.fillStyle = NODE_COLORS[node.type] || NODE_COLORS[NODE_TYPES.UNKNOWN];
              ctx.beginPath();
              ctx.arc(node.x, node.y, 5, 0, 2 * Math.PI, false);
              ctx.fill();

              ctx.font = `${fontSize}px Inter, system-ui, sans-serif`;
              ctx.textAlign = 'center';
              ctx.textBaseline = 'top';
              ctx.fillStyle = 'rgba(255,255,255,0.85)';
              ctx.fillText(label, node.x, node.y + 7);
            }}
            nodePointerAreaPaint={(node, color, ctx) => {
              ctx.fillStyle = color;
              ctx.beginPath();
              ctx.arc(node.x, node.y, 10, 0, 2 * Math.PI, false);
              ctx.fill();
            }}
            nodeLabel={(node) => {
              const vendor = node.vendor || 'Unknown vendor';
              const ip = node.ip || node.id;
              const protocols = formatProtocols(node.protocols);
              const typeLabel = NODE_LABELS[node.type] || 'Unknown';
              return `${ip}\n${vendor}\n${typeLabel}\nProtocols: ${protocols}`;
            }}
          />
        </main>
      </div>
    </div>
  );
}

export default App;
