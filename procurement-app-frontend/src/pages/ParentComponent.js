// ParentComponent.jsx
import React from 'react';
import PurchaseOrderForm from './PurchaseOrderForm'; // Gantikan dengan path yang sesuai

const ParentComponent = () => {
  const handleLogout = () => {
    // Logika logout di sini
  };

  return (
    <div>
      {/* Menggunakan PurchaseOrderForm dan memberikan prop onLogout */}
      <PurchaseOrderForm onLogout={handleLogout} />
    </div>
  );
};

export default ParentComponent;
