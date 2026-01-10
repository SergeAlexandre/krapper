import { Toolbar } from "primereact/toolbar";


export function TopBanner() {
    const startContent = (<div>LEFT</div>);
    const centerContent = (<div>CENTER</div>);
    const endContent = (<div>RIGHT</div>);

    return (
        <Toolbar start={startContent} center={centerContent} end={endContent}></Toolbar>
    )
}