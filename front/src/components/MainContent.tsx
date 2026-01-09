import type { SelectedItem } from '../types';

interface MainContentProps {
  selectedItem: SelectedItem | null;
}

export function MainContent({ selectedItem }: MainContentProps) {
  return (
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
  );
}
