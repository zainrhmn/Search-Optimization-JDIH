import { useSearchParams } from "react-router-dom";
import { useState, useEffect } from "react";
import SearchBar from "../components/SearchBar";

export default function Search() {
  const [params, setParams] = useSearchParams();

  const initialQuery = params.get("q") || "";

  const [query, setQuery] = useState(initialQuery);
  const [searchedQuery, setSearchedQuery] = useState(initialQuery);
  const [results, setResults] = useState([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [highlightOnly, setHighlightOnly] = useState(true);

  // ===== Fetch API =====
  const fetchData = async (keyword) => {
    if (!keyword.trim()) return;

    setLoading(true);

    try {
      const res = await fetch(
        `/api/regulasi/search?keyword=${encodeURIComponent(keyword)}`
      );

      const json = await res.json();

      setResults(Array.isArray(json.data) ? json.data : []);
      setTotal(json.total || 0);

      // optional: sync backend normalized query
      if (json.query) {
        setQuery(json.query);
        setSearchedQuery(json.query || keyword);
      }
    } catch (err) {
      console.error("Fetch error:", err);
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  // Load on first open / URL change
  useEffect(() => {
    fetchData(initialQuery);
  }, [initialQuery]);

  // Submit Search
  const handleSearch = (e) => {
    e.preventDefault();

    setParams({ q: query });
    fetchData(query);
  };

  return (
    <div
      className="
        min-h-screen
        bg-[linear-gradient(to_bottom,var(--color-kpu-dark),var(--color-kpu-maroon),var(--color-kpu-red))]
      "
    >
      {/* HEADER */}
      <div className="sticky top-0 z-20 backdrop-blur-md bg-black/20 border-b border-white/10">
        <div className="max-w-6xl mx-auto px-4 py-4 flex items-center gap-4">

          <img src="/logo.png" className="w-25" />

          <div className="flex-1">
            <SearchBar
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onSubmit={handleSearch}
            />
          </div>

          <label className="flex items-center gap-2 whitespace-nowrap bg-white/90 text-gray-800 px-4 py-2 rounded-full shadow-md text-sm">
            <input
              type="checkbox"
              checked={highlightOnly}
              onChange={(e) => setHighlightOnly(e.target.checked)}
              className="w-4 h-4 accent-[var(--color-kpu-gold)]"
            />
            Highlight Only
          </label>

        </div>
      </div>

      {/* BODY */}
      <div className="max-w-6xl mx-auto px-4 py-8">

        {/* Loading */}
        {loading && (
          <p className="text-white text-center">Loading...</p>
        )}

        {/* Summary */}
        {!loading && searchedQuery && (
          <p className="text-white mb-6">
            <span className="font-semibold text-[var(--color-kpu-lightGold)]">
              {total}
            </span>{" "}
            regulasi ditemukan untuk{" "}
            <span className="font-semibold text-[var(--color-kpu-lightGold)]">
              "{searchedQuery}"
            </span>
          </p>
        )}

        {/* Empty */}
        {!loading && total === 0 && (
          <p className="text-white">Tidak ada hasil ditemukan.</p>
        )}


        {/* Results */}
        <div className="space-y-8">
          {results.map((reg) => (
            <div
              key={reg.mongo_id}
              className="bg-white rounded-2xl shadow-xl p-8"
            >
              {/* Title */}
              <h1 className="text-2xl font-bold text-[var(--color-kpu-red)] mb-8">
                {reg.title}
              </h1>

              {/* Babs */}
              {reg.babs?.map((bab, i) => (
                <div key={i} className="mb-10">

                  <h2 className="text-xl font-semibold text-gray-800 mb-5 border-b pb-2">
                    BAB {bab.number} — {bab.title}
                  </h2>

                  {/* Pasals */}
                  {bab.pasals?.map((pasal, j) => {
                    const ayatsToShow = highlightOnly
                      ? pasal.ayats?.filter((ayat) => ayat.highlight)
                      : pasal.ayats;

                    if (!ayatsToShow || ayatsToShow.length === 0)
                      return null;

                    return (
                      <div
                        key={j}
                        className="mb-8 border-l-4 border-[var(--color-kpu-gold)] pl-5"
                      >
                        <h3 className="text-lg font-bold text-gray-800 mb-3">
                          Pasal {pasal.nomor}
                        </h3>

                        <div className="space-y-2">
                          {ayatsToShow.map((ayat, k) => (
                            <p
                              key={k}
                              className={
                                ayat.highlight
                                  ? "font-semibold text-black"
                                  : "text-gray-600"
                              }
                            >
                              {ayat.isi}
                            </p>
                          ))}
                        </div>
                      </div>
                    );
                  })}

                </div>
              ))}
            </div>
          ))}
        </div>

      </div>
    </div>
  );
}