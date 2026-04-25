import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import SearchBar from "../components/SearchBar";

export default function Home() {
    const [query, setQuery] = useState("");
    const navigate = useNavigate();
    const inputRef = useRef(null);

    const handleSearch = (e) => {
        e.preventDefault();
        console.log("Search:", query);

        if (!query.trim()) return;

        navigate(`/search?q=${encodeURIComponent(query)}`);
    };

    // "/" shortcut
    useEffect(() => {
        const handler = (e) => {
            if (e.key === "/") {
                e.preventDefault();
                inputRef.current?.focus();
            }
        };
        window.addEventListener("keydown", handler);
        return () => window.removeEventListener("keydown", handler);
    }, []);

    return (
        <div
            className="
        relative min-h-screen flex items-center justify-center px-4 overflow-hidden
        bg-[linear-gradient(to_bottom,var(--color-kpu-dark),var(--color-kpu-maroon),var(--color-kpu-red))]
      "
        >
            {/* Glow */}
            <div className="absolute w-[600px] h-[600px] bg-[var(--color-kpu-orange)] opacity-20 blur-3xl top-[-150px] left-1/2 -translate-x-1/2" />
            <div className="absolute w-[400px] h-[400px] bg-[var(--color-kpu-gold)] opacity-20 blur-3xl bottom-[-100px] right-[-100px]" />

            <div className="relative z-10 flex flex-col items-center text-center">

                {/* === LOGO === */}
                <img
                    src="/logo.png"
                    alt="PERISAI Logo"
                    className="w-56 drop-shadow-xl"
                />

                {/* === TITLE === */}
                <h1
                    className="
            text-6xl md:text-7xl font-extrabold tracking-tight
            bg-gradient-to-r
            from-[var(--color-kpu-gold)]
            to-[var(--color-kpu-lightGold)]
            bg-clip-text text-transparent
          "
                >
                    PERISAI
                </h1>

                {/* Divider */}
                <div
                    className="
            w-24 h-1 rounded-full mt-4
            bg-[linear-gradient(to_right,var(--color-kpu-gold),var(--color-kpu-lightGold))]
          "
                />

                {/* === UPDATED SUBTITLE === */}
                <p className="mt-5 text-gray-200 max-w-lg text-lg">
                    Pencarian Regulasi berbasis Full-Text Search untuk Akses Informasi Hukum
                </p>

                {/* Search */}
                <div className="mt-5 w-full flex justify-center">
                    <SearchBar
                        inputRef={inputRef}
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                        onSubmit={handleSearch}
                    />
                </div>

                {/* Hint */}
                <p className="mt-6 text-sm text-gray-300">
                    Tekan{" "}
                    <span className="font-mono bg-white/10 px-2 py-1 rounded">/</span>{" "}
                    untuk mulai mencari
                </p>

                {/* Example */}
                <p className="mt-2 text-sm text-gray-300">
                    Contoh:{" "}
                    <span className="italic text-[var(--color-kpu-lightGold)]">
                        "pasal 5 pkpu 17 2024"
                    </span>
                </p>

            </div>
        </div>
    );
}