import { useEffect, useState, useMemo } from 'react';
import { PrimeReactProvider } from 'primereact/api';

import { Header, Sidebar, MainContent } from './components';
import type { Catalog, SelectedItem } from './types';

import 'primereact/resources/themes/lara-dark-teal/theme.css';
import 'primeicons/primeicons.css';
import './App.css';

function App() {
  const [catalog, setCatalog] = useState<Catalog | null>(null);
  const [selectedItem, setSelectedItem] = useState<SelectedItem | null>(null);

  useEffect(() => {
    fetch('/api/v1/wraps')
      .then(res => res.json())
      .then((data: Catalog) => setCatalog(data))
      .catch(err => console.error('Failed to fetch catalog:', err));
  }, []);

  const wraps = useMemo(() => catalog?.wraps ?? [], [catalog]);

  return (
    <PrimeReactProvider>
      <div className="app-container">
        <Header />
        <div className="app-body">
          <Sidebar
            catalog={wraps}
            onSelectItem={setSelectedItem}
          />
          <MainContent selectedItem={selectedItem} />
        </div>
      </div>
    </PrimeReactProvider>
  );
}

export default App;
