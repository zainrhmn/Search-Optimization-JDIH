export default function SearchBar({
  value,
  onChange,
  onSubmit,
  inputRef,
}) {
  return (
    <form onSubmit={onSubmit} className="w-full max-w-2xl">
      <div className="relative group">

        <input
          ref={inputRef}
          type="text"
          value={value}
          onChange={onChange}
          placeholder="Cari pasal, ayat, atau kata kunci..."
          className="
            w-full px-6 py-5 pr-14
            rounded-full
            bg-white/80 backdrop-blur-md
            border border-gray-200
            shadow-md

            text-gray-800 placeholder-gray-400

            focus:outline-none
            focus:ring-2 focus:ring-blue-500
            focus:border-transparent

            transition-all duration-200
            group-hover:shadow-lg
          "
        />

        {/* search icon */}
        <span className="absolute right-5 top-1/2 -translate-y-1/2 text-gray-400 group-focus-within:text-blue-500 transition">
          🔍
        </span>

      </div>
    </form>
  );
}