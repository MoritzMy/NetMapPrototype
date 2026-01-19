import {
  memo,
  useCallback,
  useEffect,
  useLayoutEffect,
  useMemo,
  useRef,
  useState
} from 'react';
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

const GraphView = memo(function GraphView({ graphData, dimensions, nodeRelSize, graphRef }) {
  const nodeRadius = Math.max(4, nodeRelSize);

  return (
    <ForceGraph2D
      ref={graphRef}
      graphData={graphData}
      nodeId="id"
      linkSource="source"
      linkTarget="target"
      width={dimensions.width || 800}
      height={dimensions.height || 600}
      backgroundColor="#0b0f1c"
      nodeRelSize={nodeRelSize}
      linkDirectionalParticles={2}
      linkDirectionalParticleWidth={2}
      linkDirectionalParticleSpeed={0.004}
      linkColor={(link) => linkColorFor(link)}
      linkWidth={(link) => (link.type === 'routes-via' ? 2 : 1.25)}
      nodeCanvasObject={(node, ctx, globalScale) => {
        const label = node.ip || node.id;
        const fontSize = 12 / globalScale;
        const radius = nodeRadius / globalScale;

        ctx.fillStyle = NODE_COLORS[node.type] || NODE_COLORS[NODE_TYPES.UNKNOWN];
        ctx.beginPath();
        ctx.arc(node.x, node.y, radius, 0, 2 * Math.PI, false);
        ctx.fill();

        ctx.font = `${fontSize}px Inter, system-ui, sans-serif`;
        ctx.textAlign = 'center';
        ctx.textBaseline = 'top';
        ctx.fillStyle = 'rgba(255,255,255,0.85)';
        ctx.fillText(label, node.x, node.y + radius + 2);
      }}
      nodePointerAreaPaint={(node, color, ctx, globalScale) => {
        const radius = (nodeRadius + 4) / globalScale;
        ctx.fillStyle = color;
        ctx.beginPath();
        ctx.arc(node.x, node.y, radius, 0, 2 * Math.PI, false);
        ctx.fill();
      }}
      nodeLabel={(node) => {
        const vendor = node.vendor || 'Unknown vendor';
        const ip = node.ip || node.id;
        const mac = node.mac || 'Unknown MAC';
        const protocols = formatProtocols(node.protocols);
        const typeLabel = NODE_LABELS[node.type] || 'Unknown';
        return `${ip}\n${mac}\n${vendor}\n${typeLabel}\nProtocols: ${protocols}`;
      }}
    />
  );
});

function App() {
  const [graphData, setGraphData] = useState(DEFAULT_GRAPH);
  const [status, setStatus] = useState('idle');
  const [lastUpdated, setLastUpdated] = useState(null);
  const [scanStatus, setScanStatus] = useState('');
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
  const [forceRefreshToken, setForceRefreshToken] = useState(0);
  const graphRef = useRef(null);
  const graphContainerRef = useRef(null);
  const nodeIdSetRef = useRef(new Set());

  const fetchGraph = useCallback(async ({ force = false } = {}) => {
    setStatus('loading');
    try {
      const response = await fetch('/api/graph');
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }
      const data = await response.json();
      const nextGraph = data ?? DEFAULT_GRAPH;
      const nextNodes = nextGraph.nodes ?? [];
      const nextIds = new Set(nextNodes.map((node) => node.id));
      const existingIds = nodeIdSetRef.current;
      const hasNewNode = Array.from(nextIds).some((id) => !existingIds.has(id));

      if (force || hasNewNode || existingIds.size === 0) {
        nodeIdSetRef.current = nextIds;
        setGraphData(nextGraph);
      }
      setLastUpdated(new Date());
      setStatus('success');
    } catch (error) {
      setStatus('error');
    }
  }, []);

  const triggerScan = useCallback(
    async (endpoint, label) => {
      setScanStatus(`${label} requested...`);
      try {
        const response = await fetch(endpoint, { method: 'POST' });
        if (!response.ok) {
          throw new Error(`HTTP ${response.status}`);
        }
        setScanStatus(`${label} started.`);
        fetchGraph();
      } catch (error) {
        setScanStatus(`${label} failed. Check the backend service.`);
      }
    },
    [fetchGraph]
  );

  useEffect(() => {
    fetchGraph({ force: true });
    const interval = setInterval(() => fetchGraph({ force: false }), 10000);
    return () => clearInterval(interval);
  }, [fetchGraph]);

  useEffect(() => {
    if (graphRef.current && graphData.nodes.length > 0) {
      requestAnimationFrame(() => {
        graphRef.current.centerAt(0, 0, 400);
        graphRef.current.zoomToFit(500, 120);
      });
    }
  }, [graphData, forceRefreshToken]);

  useLayoutEffect(() => {
    const container = graphContainerRef.current;
    if (!container) return;

    const observer = new ResizeObserver((entries) => {
      if (!entries.length) return;
      const { width, height } = entries[0].contentRect;
      setDimensions({
        width: Math.max(1, Math.floor(width)),
        height: Math.max(1, Math.floor(height))
      });
    });

    observer.observe(container);
    return () => observer.disconnect();
  }, []);

  const nodeRelSize = useMemo(() => {
    const count = graphData.nodes.length;
    if (count <= 10) return 12;
    if (count <= 50) return 10;
    if (count <= 150) return 8;
    if (count <= 300) return 6;
    if (count <= 600) return 5;
    return 4;
  }, [graphData.nodes.length]);

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
          <button
            type="button"
            onClick={() => {
              setForceRefreshToken((value) => value + 1);
              fetchGraph({ force: true });
            }}
            disabled={status === 'loading'}
          >
            {status === 'loading' ? 'Refreshing…' : 'Refresh'}
          </button>
          <button
            className="secondary"
            type="button"
            onClick={() => triggerScan('/api/icmp-sweep', 'ICMP sweep')}
          >
            Run ICMP Sweep
          </button>
          <button
            className="secondary"
            type="button"
            onClick={() => triggerScan('/api/arp-scan', 'ARP scan')}
          >
            Run ARP Scan
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
      {scanStatus ? <div className="scan-status">{scanStatus}</div> : null}

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

        <main className="graph-panel" ref={graphContainerRef}>
          <GraphView
            graphData={graphData}
            dimensions={dimensions}
            nodeRelSize={nodeRelSize}
            graphRef={graphRef}
          />
        </main>
      </div>
    </div>
  );
}

export default App;
