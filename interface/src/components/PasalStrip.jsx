// src/components/PasalStrip.jsx

import React from "react";

const PasalStrip = ({ bab, pasals, startIndex, setStartIndex, selectedPasal, onSelectPasal }) => {
  const visiblePasals = pasals.slice(startIndex, startIndex + 10);
  const isBab3 = bab.number === "3";

  return (
    <div className="pasal-strip-wrapper">
      {isBab3 && (
        <button
          className="arrow-btn"
          onClick={() => setStartIndex((prev) => Math.max(prev - 10, 0))}
          disabled={startIndex === 0}
        >
          ◀
        </button>
      )}

      <div className="pasal-strip">
        {visiblePasals.map((pasal) => (
          <div
            key={pasal.number}
            className={`pasal-box ${selectedPasal?.number === pasal.number ? "active" : ""}`}
            onClick={() => onSelectPasal(pasal)}
          >
            Pasal {pasal.number}
          </div>
        ))}
      </div>

      {isBab3 && startIndex + 10 < pasals.length && (
        <button
          className="arrow-btn"
          onClick={() => setStartIndex((prev) => Math.min(prev + 10, pasals.length - 10))}
        >
          ▶
        </button>
      )}
    </div>
  );
};

export default PasalStrip;
