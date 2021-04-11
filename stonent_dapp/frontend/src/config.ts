enum MODE {
    DEV = 'DEV',
    TEST = 'TEST',
    PROD = 'PROD',
}

function getConfig(mode: MODE) {
    switch (mode) {
        case MODE.DEV: {
            return {
                mode,
                BACKEND_URL: '',
            };
        }
        case MODE.TEST: {
            return {
                mode,
                BACKEND_URL: '',
            };
        }
        case MODE.PROD: {
            return {
                mode,
                BACKEND_URL: '',
            };
        }
    }
}

export default getConfig(MODE.DEV);
