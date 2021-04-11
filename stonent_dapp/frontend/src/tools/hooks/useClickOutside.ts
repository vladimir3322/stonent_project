import {RefObject, useEffect} from 'react';


export default function<IRef extends HTMLElement>(ref: RefObject<IRef>, handler: () => void) {
    const handleClick = (e: any) => {
        if (ref.current && !ref.current.contains(e.target)) {
            handler();
        }
    };

    useEffect(() => {
        document.addEventListener('click', handleClick);

        return () => {
            document.removeEventListener('click', handleClick);
        };
    });
}
