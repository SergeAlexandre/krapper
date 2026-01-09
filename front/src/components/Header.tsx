import { Avatar } from 'primereact/avatar';
import { Button } from 'primereact/button';

export function Header() {
  return (
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
  );
}
