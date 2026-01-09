import { useEffect, useState } from 'react';
import { PrimeReactProvider } from 'primereact/api';
import { PanelMenu } from 'primereact/panelmenu';
import { Avatar } from 'primereact/avatar';
import { Button } from 'primereact/button';

import 'primereact/resources/themes/lara-dark-teal/theme.css';
import 'primeicons/primeicons.css';
import './App.css';

interface MenuItem {
  key?: string;
  label?: string;
  icon?: string;
  command?: () => void;
  items?: MenuItem[];
  disabled?: boolean;
  expanded?: boolean;
}

interface CatalogItem {
  name: string;
  label: string;
  menuMode: 'grid' | 'subMenu';
}

interface Catalog {
  wraps: CatalogItem[];
}

interface K8sResource {
  metadata: {
    name: string;
    namespace?: string;
  };
}

type SubMenuState = {
  loading: boolean;
  error?: string;
  resources: K8sResource[];
};

function App() {
  const [catalog, setCatalog] = useState<Catalog | null>(null);
  const [subMenuState, setSubMenuState] = useState<Record<string, SubMenuState>>({});
  const [expandedKeys, setExpandedKeys] = useState<Record<string, boolean>>({});
  const [selectedItem, setSelectedItem] = useState<{ wrap: string; resource?: string } | null>(null);

  useEffect(() => {
    fetch('/api/v1/wraps')
      .then(res => res.json())
      .then((data: Catalog) => setCatalog(data))
      .catch(err => console.error('Failed to fetch catalog:', err));
  }, []);

  const fetchSubMenuResources = async (wrapName: string) => {
    const state = subMenuState[wrapName];
    if (state && (state.loading || state.resources.length > 0)) return; // Already fetching or fetched
    
    // Set loading state
    setSubMenuState(prev => ({
      ...prev,
      [wrapName]: { loading: true, resources: [] }
    }));
    
    try {
      const res = await fetch(`/api/v1/resources/${wrapName}`);
      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(errorText || `HTTP ${res.status}`);
      }
      const resources: K8sResource[] = await res.json();
      setSubMenuState(prev => ({
        ...prev,
        [wrapName]: { loading: false, resources }
      }));
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Unknown error';
      console.error(`Failed to fetch resources for ${wrapName}:`, errorMsg);
      setSubMenuState(prev => ({
        ...prev,
        [wrapName]: { loading: false, error: errorMsg, resources: [] }
      }));
    }
  };

  const menuItems: MenuItem[] = catalog?.wraps.map(item => {
    if (item.menuMode === 'subMenu') {
      const state = subMenuState[item.name];
      const isLoading = state?.loading ?? false;
      const hasError = !!state?.error;
      const resources = state?.resources ?? [];
      const hasResources = resources.length > 0;
      
      let subItems: MenuItem[];
      if (isLoading) {
        subItems = [{ label: 'Loading...', icon: 'pi pi-spin pi-spinner', disabled: true }];
      } else if (hasError) {
        subItems = [{ label: 'Error loading resources', icon: 'pi pi-exclamation-triangle', disabled: true }];
      } else if (hasResources) {
        subItems = resources.map(resource => ({
          key: `${item.name}/${resource.metadata.name}`,
          label: resource.metadata.name,
          icon: 'pi pi-file',
          command: () => setSelectedItem({ wrap: item.name, resource: resource.metadata.name })
        }));
      } else {
        subItems = [{ label: 'No resources found', icon: 'pi pi-info-circle', disabled: true }];
      }
      
      return {
        key: item.name,
        label: item.label,
        icon: 'pi pi-folder',
        expanded: expandedKeys[item.name],
        command: () => {
          // Toggle expanded state and fetch if needed
          const isExpanding = !expandedKeys[item.name];
          setExpandedKeys(prev => ({ ...prev, [item.name]: isExpanding }));
          if (isExpanding) {
            fetchSubMenuResources(item.name);
          }
        },
        items: subItems
      };
    }
    
    // grid mode - simple menu item
    return {
      key: item.name,
      label: item.label,
      icon: 'pi pi-th-large',
      command: () => setSelectedItem({ wrap: item.name })
    };
  }) ?? [];

  return (
    <PrimeReactProvider>
      <div className="app-container">
        {/* Top Banner */}
        <header className="app-header">
          <div className="header-left">
            <i className="pi pi-box logo-icon"></i>
            <span className="logo-text">Krapper</span>
          </div>
          <div className="header-right">
            <Button icon="pi pi-cog" rounded text aria-label="Settings" />
            <Avatar icon="pi pi-user" shape="circle" />
          </div>
        </header>

        <div className="app-body">
          {/* Left Sidebar */}
          <aside className="app-sidebar">
            <PanelMenu 
              model={menuItems} 
              className="sidebar-menu"
              multiple
            />
          </aside>

          {/* Main Content */}
          <main className="app-content">
            {selectedItem ? (
              <div className="content-placeholder">
                <h2>{selectedItem.wrap}</h2>
                {selectedItem.resource && <p>Resource: {selectedItem.resource}</p>}
              </div>
            ) : (
              <div className="content-placeholder">
                <p>Select an item from the menu</p>
              </div>
            )}
          </main>
        </div>
      </div>
    </PrimeReactProvider>
  );
}

export default App;
