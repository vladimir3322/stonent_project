import {useEffect, useRef} from 'react';


export default function<IValue>(value: IValue): IValue | null {
    const ref = useRef<IValue>();

    useEffect(() => {
        ref.current = value;
    });

    return ref.current || null;
}
