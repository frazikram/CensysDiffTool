import { useState, useEffect } from "react";
import { apiClient, Snapshot, DiffResult } from "../api/client";

export default function HistoryPage() {
  const [hosts, setHosts] = useState<string[]>([]);
  const [ip, setIp] = useState("");
  const [snapshots, setSnapshots] = useState<Snapshot[]>([]);
  const [selected, setSelected] = useState<string[]>([]);
  const [diff, setDiff] = useState<DiffResult | null>(null);
  const [hasSearched, setHasSearched] = useState(false); // NEW

  // Fetch list of all hosts when page loads
  useEffect(() => {
    apiClient
      .getHosts()
      .then(setHosts)
      .catch((err) => console.error("Failed to fetch hosts", err));
  }, []);

  const fetchHistory = async () => {
    if (!ip) return;
    setSnapshots([]);
    setSelected([]);
    setDiff(null);
    setHasSearched(true); // mark that a search happened
    try {
      const result = await apiClient.getHistory(ip);
      setSnapshots(result);
    } catch (err: any) {
      alert(`Failed to fetch history: ${err.message}`);
    }
  };

  const toggleSelect = (id: string) => {
    setSelected(
      (prev) =>
        prev.includes(id)
          ? prev.filter((x) => x !== id)
          : [...prev, id].slice(-2) // max 2
    );
  };

  const generateDiff = async () => {
    if (selected.length === 2) {
      try {
        const result = await apiClient.getDiff(selected[0], selected[1]);
        setDiff(result);
      } catch (err: any) {
        alert(`Failed to generate diff: ${err.message}`);
      }
    }
  };

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <h1 className="text-xl font-bold mb-4">Host History</h1>

      {/* Host Selector */}
      <div className="flex gap-2 mb-4">
        <select
          value={ip}
          onChange={(e) => setIp(e.target.value)}
          className="border px-2 py-1 rounded flex-1"
        >
          <option value="">-- Select Host --</option>
          {hosts.map((h) => (
            <option key={h} value={h}>
              {h}
            </option>
          ))}
        </select>
        <input
          type="text"
          value={ip}
          onChange={(e) => setIp(e.target.value)}
          placeholder="Or enter manually"
          className="border px-2 py-1 rounded flex-1"
        />
        <button
          onClick={fetchHistory}
          className="bg-blue-600 text-white px-4 py-2 rounded"
        >
          Search
        </button>
      </div>

      {/* Snapshots */}
      {snapshots.length > 0 && (
        <ul className="space-y-2 mb-4">
          {snapshots.map((s) => (
            <li
              key={s.id}
              className="p-2 border rounded flex items-center justify-between"
            >
              <span>{new Date(s.timestamp).toLocaleString()}</span>
              <input
                type="checkbox"
                checked={selected.includes(s.id!)}
                onChange={() => toggleSelect(s.id!)}
              />
            </li>
          ))}
        </ul>
      )}

      {/* Only show if user searched and no results */}
      {hasSearched && snapshots.length === 0 && (
        <p>No snapshots found for {ip}</p>
      )}

      {/* Diff Button */}
      <button
        disabled={selected.length !== 2}
        onClick={generateDiff}
        className="bg-green-600 text-white px-4 py-2 rounded disabled:opacity-50"
      >
        Generate Diff Report
      </button>

      {diff && (
  <div className="mt-6 border-t pt-4">
    <h2 className="text-lg font-bold mb-2">Diff Report</h2>
    <div>
      <h3 className="text-green-600 font-semibold">Added</h3>
      {(diff.added ?? []).length > 0 ? (
        (diff.added ?? []).map((s, i) => (
          <p key={i}>
            {s.protocol} on port {s.port}
          </p>
        ))
      ) : (
        <p className="text-gray-500">None</p>
      )}

      <h3 className="text-red-600 font-semibold mt-2">Removed</h3>
      {(diff.removed ?? []).length > 0 ? (
        (diff.removed ?? []).map((s, i) => (
          <p key={i}>
            {s.protocol} on port {s.port}
          </p>
        ))
      ) : (
        <p className="text-gray-500">None</p>
      )}

      <h3 className="text-yellow-600 font-semibold mt-2">Changed</h3>
      {(diff.changed ?? []).length > 0 ? (
        (diff.changed ?? []).map((c, i) => (
          <div key={i}>
            <p>
              Before: {c.from.protocol} {c.from.software?.version}
            </p>
            <p>
              After: {c.to.protocol} {c.to.software?.version}
            </p>
          </div>
        ))
      ) : (
        <p className="text-gray-500">None</p>
      )}
    </div>
  </div>
)}
  </div>
  );
}
