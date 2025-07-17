import React, { useState } from "react";
import regulationData from "./data/regulationData";
import "./style.css";

const App = () => {
  const [selectedBabIndex, setSelectedBabIndex] = useState(0);
  const [selectedPasalIndex, setSelectedPasalIndex] = useState(null);
  const [pasalStartIndex, setPasalStartIndex] = useState(0);
  const PASAL_LIMIT = 10;

  const handleBabChange = (e) => {
    const index = parseInt(e.target.value);
    setSelectedBabIndex(index);
    setSelectedPasalIndex(null);
    setPasalStartIndex(0);
  };

  const handlePasalClick = (index) => {
    setSelectedPasalIndex(index);
  };

  const handleNext = () => {
    if (
      pasalStartIndex + PASAL_LIMIT <
      regulationData.babs[selectedBabIndex].pasals.length
    ) {
      setPasalStartIndex(pasalStartIndex + PASAL_LIMIT);
    }
  };

  const handlePrev = () => {
    if (pasalStartIndex > 0) {
      setPasalStartIndex(pasalStartIndex - PASAL_LIMIT);
    }
  };

  const bab = regulationData.babs[selectedBabIndex];
  const visiblePasals = bab.number === 3
    ? bab.pasals.slice(pasalStartIndex, pasalStartIndex + PASAL_LIMIT)
    : bab.pasals;

  return (
    <div className="app-container">
      <header className="header">
        <img src="/logo kpu.png" alt="Logo KPU" className="logo" />
        <h1 className="header-title">
          Peraturan Komisi Pemilihan Umum Nomor 1 Tahun 2025
        </h1>
      </header>

      <div className="dropdown-container">
        <label htmlFor="bab-select">Pilih BAB:</label>
        <select id="bab-select" onChange={handleBabChange} value={selectedBabIndex}>
          {regulationData.babs.map((bab, index) => (
            <option key={index} value={index}>
              BAB {bab.number}: {bab.title}
            </option>
          ))}
        </select>
      </div>

      <div className="pasal-strip-wrapper">
        {bab.number === 3 && (
          <button onClick={handlePrev} className="arrow-btn">◀</button>
        )}

        <div className="pasal-strip">
          {visiblePasals.map((pasal, index) => {
            const realIndex = bab.number === 3 ? index + pasalStartIndex : index;
            return (
              <button
                key={realIndex}
                className={`pasal-box ${realIndex === selectedPasalIndex ? "active" : ""}`}
                onClick={() => handlePasalClick(realIndex)}
              >
                {pasal.title}
              </button>
            );
          })}
        </div>

        {bab.number === 3 && (
          <button onClick={handleNext} className="arrow-btn">▶</button>
        )}
      </div>

      <div className="content-box">
        {selectedPasalIndex !== null && (
          <>
            <h2>{bab.pasals[selectedPasalIndex].title}</h2>
            <p>{bab.pasals[selectedPasalIndex].description}</p>
            <ul>
              {bab.pasals[selectedPasalIndex].ayats.map((ayat, i) => (
                <li key={i}>{ayat}</li>
              ))}
            </ul>
          </>
        )}
      </div>
    </div>
  );
};

export default App;
