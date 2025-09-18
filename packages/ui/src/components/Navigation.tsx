import { Link } from "react-router-dom";

export default function Navigation() {
  return (
    <nav className="bg-gray-900 text-white p-4 flex gap-4">
      <Link to="/upload" className="hover:text-blue-400">Upload</Link>
      <Link to="/history/203.0.113.45" className="hover:text-blue-400">Sample History</Link>
    </nav>
  );
}
