
import { PanelMenu } from 'primereact/panelmenu';
import type { MenuItem } from 'primereact/menuitem';
import {ScrollPanel} from "primereact/scrollpanel";

export default function LeftMenu() {
    let items: MenuItem[] = [
        {
            label: 'Files',
            icon: 'pi pi-file',
            className: 'submenu-0',
            items: [
                {
                    label: 'Documents',
                    icon: 'pi pi-file',
                    className: 'submenu-1',
                    items: [
                        {
                            label: 'Invoices',
                            icon: 'pi pi-file-pdf',
                            className: 'submenu-2',
                            items: [
                                {
                                    label: 'Pending',
                                    icon: 'pi pi-stop',
                                    className: 'submenu-3',
                                },
                                {
                                    label: 'Paid',
                                    icon: 'pi pi-check-circle',
                                    className: 'submenu-3',
                                }
                            ]
                        },
                        {
                            label: 'Clients',
                            className: 'submenu-2',
                            icon: 'pi pi-users'
                        }
                    ]
                },
                {
                    label: 'Images',
                    icon: 'pi pi-image',
                    className: 'submenu-1',
                    items: [
                        {
                            label: 'Logos',
                            className: 'submenu-2',
                            icon: 'pi pi-image'
                        }
                    ]
                }
            ]
        },
        {
            label: 'Cloud',
            icon: 'pi pi-cloud',
            className: 'submenu-0',
            items: [
                {
                    label: 'Upload',
                    icon: 'pi pi-cloud-upload',
                    className: 'submenu-1',
                },
                {
                    label: 'Download',
                    icon: 'pi pi-cloud-download',
                    className: 'submenu-1',
                },
                {
                    label: 'Sync',
                    icon: 'pi pi-refresh',
                    className: 'submenu-1',
                }
            ]
        },
        {
            label: 'Devices',
            icon: 'pi pi-desktop',
            className: 'submenu-0',
            items: [
                {
                    label: 'Phone',
                    icon: 'pi pi-mobile',
                    className: 'submenu-1',
                },
                {
                    label: 'Desktop',
                    icon: 'pi pi-desktop',
                    className: 'submenu-1',
                },
                {
                    label: 'Tablet',
                    icon: 'pi pi-tablet',
                    className: 'submenu-1',
                }
            ]
        },
        {
            label: 'Direct',
            icon: 'pi pi-desktop',
        }
    ];
    return (
        <div className="card flex justify-content-center" style={{ height: "100%" }}>
            <ScrollPanel style={{ width: '20rem', height: '100%' }} >
                <div>
                    <PanelMenu model={items} multiple />
                </div>
            </ScrollPanel>
        </div>
    )
}
