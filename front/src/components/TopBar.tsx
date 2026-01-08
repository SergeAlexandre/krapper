
import { Menubar } from 'primereact/menubar';
import { Avatar } from 'primereact/avatar';

const TopBar = () => {
    const start = <div className="font-bold text-xl">My App</div>;

    const end = (
        <div className="flex align-items-center gap-2">
            <i className="pi pi-cog p-element" style={{ fontSize: '1.5rem', cursor: 'pointer', marginRight: '1rem' }} />
            <div className="flex align-items-center gap-2 cursor-pointer">
                <span>My Account</span>
                <Avatar icon="pi pi-user" shape="circle" />
            </div>
        </div>
    );

    return (
        <div className="card">
            <Menubar start={start} end={end} style={{ borderRadius: 0, border: 'none', borderBottom: '1px solid #dee2e6' }} />
        </div>
    );
}

export default TopBar;
