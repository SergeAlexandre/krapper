export interface CatalogItem {
  name: string;
  label: string;
  menuMode: 'grid' | 'subMenu';
}

export interface Catalog {
  wraps: CatalogItem[];
}

export interface SelectedItem {
  wrap: string;
  resource?: string;
}
