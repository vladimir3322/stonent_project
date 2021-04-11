import {useEffect, useRef} from 'react';


type IAnyFunction = <IArgs extends unknown[], IResult = unknown>(...args: IArgs) => IResult;

export default function(callback: IAnyFunction, delay: number) {
    const savedCallback = useRef<IAnyFunction>();

    useEffect(() => {
        savedCallback.current = callback;
    }, [callback]);

    useEffect(() => {
        function tick() {
            savedCallback.current?.();
        }

        if (delay !== null) {
            const id = setInterval(tick, delay);

            return () => clearInterval(id);
        }
    }, [delay]);
}
