// components/RegulationViewer.jsx

import React from 'react';

const RegulationViewer = ({ data }) => {
  return (
    <div className="container">
      <h1>{data.title}</h1>
      <p className="meta-info">Status: {data.status} • Tahun: {data.tahun}</p>

      {data.babs.map((bab, indexBab) => (
        <div key={indexBab}>
          <h2 className="bab-title">BAB {bab.number}: {bab.title}</h2>

          {bab.pasals.map((pasal, indexPasal) => (
            <div className="pasal-card" key={indexPasal}>
              <div className="pasal-title">{pasal.title}</div>
              <div className="pasal-description">{pasal.description}</div>

              <ul className="ayat-list">
                {pasal.ayats.map((ayat, indexAyat) => (
                  <li key={indexAyat}>{ayat}</li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      ))}
    </div>
  );
};

export default RegulationViewer;
