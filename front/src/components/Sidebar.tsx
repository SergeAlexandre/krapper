
import { Menu } from 'primereact/menu';
import { useEffect, useState } from 'react';
import type { MenuItem } from 'primereact/menuitem';

interface Wrap {
    name: string;
}

const Sidebar = () => {
    const [items, setItems] = useState<MenuItem[]>([]);

    useEffect(() => {
        fetch('/api/v1/wraps')
            .then(res => res.json())
            .then((data: { wraps: Wrap[] }) => {
                const menuItems = [
                    { label: 'Dashboard', icon: 'pi pi-home' },
                    ...data.wraps.map((wrap) => ({
                        label: wrap.name,
                        icon: 'pi pi-box'
                    }))
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
