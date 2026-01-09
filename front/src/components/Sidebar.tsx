
import { Menu } from 'primereact/menu';
import { useEffect, useState } from 'react';
import type { MenuItem } from 'primereact/menuitem';

interface Wrap {
    name: string;
    menuMode?: string;
}

const Sidebar = () => {
    const [items, setItems] = useState<MenuItem[]>([]);

    useEffect(() => {
        fetch('/api/v1/wraps')
            .then(res => res.json())
            .then((data: { wraps: Wrap[] }) => {
                const processWrap = (wrap: Wrap): MenuItem => {
                    const item: MenuItem = {
                        label: wrap.name,
                        icon: 'pi pi-box'
                    };

                    if (wrap.menuMode === 'subMenu') {
                        console.log("#################1")
                        item.items = [{ label: 'Loading...', icon: 'pi pi-spin pi-spinner' }];
                        item.command = () => {
                            // Only fetch if we haven't already loaded real data (check if first item is "Loading...")
                            // Note: e.item.items might be undefined in type def but accessible at runtime,
                            // or we can track loaded state separately. simpler here: fetch always or check children.
                            // Actually, PrimeReact menu command might not easily give us access to the state update target
                            // purely from 'e'. We need to update the main 'items' state.
                            console.log("#################2")

                            fetch(`/api/v1/resources/${wrap.name}`)
                                .then(res => res.json())
                                .then((resources: { metadata: { name: string } }[]) => {
                                    setItems(prevItems => {
                                        return prevItems.map(prevItem => {
                                            if (prevItem.label === wrap.name) {
                                                return {
                                                    ...prevItem,
                                                    items: resources.map(r => ({
                                                        label: r.metadata.name,
                                                        icon: 'pi pi-file'
                                                    }))
                                                };
                                            }
                                            return prevItem;
                                        });
                                    });
                                })
                                .catch(err => console.error("Failed to load resources", err));
                        };
                    }
                    return item;
                };

                const menuItems = [
                    { label: 'Dashboard', icon: 'pi pi-home' },
                    ...data.wraps.map(processWrap)
                ];
                setItems(menuItems);
            })
            .catch(err => {
                console.error("Failed to fetch wraps:", err);
                // Fallback for demo if fetch fails
                setItems([
                    { label: 'Dashboard', icon: 'pi pi-home' },
                    { label: 'Wraps (Error)', icon: 'pi pi-exclamation-triangle' }
                ]);
            });
    }, []);

    return (
        <div className="h-full border-right-1 surface-border">
            <Menu model={items} style={{ width: '100%', border: 'none' }} />
        </div>
    );
}

export default Sidebar;
