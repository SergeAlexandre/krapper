import { useState } from 'react';
import { PanelMenu } from 'primereact/panelmenu';
import type { CatalogItem, SelectedItem } from '../types';

interface MenuItem {
  key?: string;
  label?: string;
  icon?: string;
  command?: () => void;
  items?: MenuItem[];
  disabled?: boolean;
  expanded?: boolean;
}

interface SubMenuState {
  loading: boolean;
  error?: string;
  resources: K8sResource[];
}

interface K8sResource {
  metadata: {
    name: string;
    namespace?: string;
  };
}

interface SidebarProps {
  catalog: CatalogItem[];
  onSelectItem: (item: SelectedItem) => void;
}

export function Sidebar({ catalog, onSelectItem }: SidebarProps) {
  const [subMenuState, setSubMenuState] = useState<Record<string, SubMenuState>>({});
  const [expandedKeys, setExpandedKeys] = useState<Record<string, boolean>>({});

  const fetchSubMenuResources = async (wrapName: string) => {
    const state = subMenuState[wrapName];
    if (state?.loading) return;

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

  const menuItems: MenuItem[] = catalog.map(item => {
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
          command: () => onSelectItem({ wrap: item.name, resource: resource.metadata.name })
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
          const isExpanding = !expandedKeys[item.name];
          setExpandedKeys(prev => ({ ...prev, [item.name]: isExpanding }));
          if (isExpanding) {
            fetchSubMenuResources(item.name);
          }
        },
        items: subItems
      };
    }

    return {
      key: item.name,
      label: item.label,
      icon: 'pi pi-th-large',
      command: () => onSelectItem({ wrap: item.name })
    };
  });

  return (
    <aside className="app-sidebar">
      <PanelMenu model={menuItems} className="sidebar-menu" multiple />
    </aside>
  );
}
