import {IProps} from './types';
import {IImageMetadata} from 'instances/raribleContract/types';
import {ICertificate} from 'instances/stonentContract/types';

import {Button, Container, Typography} from '@material-ui/core';
import React, {FC, useEffect, useState} from 'react';
import {useIntl} from 'react-intl';
import {useHistory} from 'react-router-dom';

import {useIercContract} from 'instances/iercContract/hooks';
import {getImageUrl} from 'instances/raribleContract/functions';
import {useRaribleContract} from 'instances/raribleContract/hooks';
import {useStonentContract} from 'instances/stonentContract/hooks';

import ProgressCertificate from 'components/Check/ProgressCertificate';
import Header from 'components/Helpers/Header';
import Spinner from 'components/Helpers/Spinner';

import Certificate from './Certificate';
import EmptyCertificate from './EmptyCertificate';
import ErrorMessage from './ErrorMessage';
import SuccessMessage from './SuccessMessage';

import {certificateHasBeenUpdated} from './functions';
import {INTL_DATA} from './intl';

import styles from './Check.module.scss';


const Check: FC<IProps> = (props) => {
    const [imageMetadata, setImageMetadata] = useState<IImageMetadata | null>(null);
    const [imageCertificate, setImageCertificate] = useState<ICertificate | null>(null);
    const [isFetching, setIsFetching] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');
    const [successMessage, setSuccessMessage] = useState('');

    const {raribleContract, getImageMetadataUrl, getImageMetadata} = useRaribleContract();
    const {iercContract, approveMoneyUsage, checkAccessMoneyUsage} = useIercContract();
    const {stonentContract, certificationPrice, getCertificate, checkImage, getPrice, subscribeToCertificateSent, subscribeToCertified} = useStonentContract();

    const history = useHistory();
    const intl = useIntl();

    const runCheckImage = async () => {
        if (!iercContract) {
            return;
        }

        setIsFetching(true);

        const alreadyApproved = await checkAccessMoneyUsage();

        if (!alreadyApproved) {
            try {
                await approveMoneyUsage();
            } catch (error) {
                setErrorMessage(intl.formatMessage({id: INTL_DATA.ERROR_PAYMENT}));
                setIsFetching(false);

                return;
            }
        }

        try {
            await checkImage(props.match.params.id);
        } catch {
            setErrorMessage(intl.formatMessage({id: INTL_DATA.ERROR_ABORTED}));
        }

        setIsFetching(false);
    };

    useEffect(() => {
        (async () => {
            if (!raribleContract) {
                return;
            }

            const metadataUrl = await getImageMetadataUrl(props.match.params.id);

            if (!metadataUrl) {
                history.push('/not_found');

                return;
            }

            const metadata = await getImageMetadata(metadataUrl);

            if (!metadata) {
                history.push('/not_found');

                return;
            }

            setImageMetadata(metadata);
        })();
    }, [raribleContract]);
    useEffect(() => {
        (async () => {
            if (!imageMetadata) {
                return;
            }
            if (!stonentContract) {
                return;
            }

            subscribeToCertificateSent(async (data) => {
                const newCertificate = await getCertificate(props.match.params.id);

                if (certificateHasBeenUpdated(imageCertificate, newCertificate)) {
                    setImageCertificate(newCertificate);
                    setSuccessMessage(intl.formatMessage({id: INTL_DATA.SUCCESS_UPDATED}));
                }
            });
            subscribeToCertified(async (data) => {
                const newCertificate = await getCertificate(props.match.params.id);

                if (certificateHasBeenUpdated(imageCertificate, newCertificate)) {
                    setImageCertificate(newCertificate);
                    setSuccessMessage(intl.formatMessage({id: INTL_DATA.SUCCESS_UPDATED}));
                }
            });

            const certificate = await getCertificate(props.match.params.id);

            if (!certificate) {
                return;
            }

            setImageCertificate(certificate);
        })();
        (async () => {
            if (!imageMetadata) {
                return;
            }
            if (!stonentContract) {
                return;
            }

            await getPrice();
        })();
    }, [imageMetadata, stonentContract]);

    if (!imageMetadata || !imageCertificate) {
        return (
            <Spinner isPage={true}/>
        );
    }

    return (
        <>
            <Header/>
            <Container>
                <div className={styles.imagePart}>
                    <div className={styles.image}>
                        <img src={getImageUrl(imageMetadata.imageIPFSURL) || undefined} alt={''}/>
                    </div>
                    <div className={styles.info}>
                        <div className={styles.top}>
                            <div>
                                <Typography variant={'h4'}>
                                    {imageMetadata.name}
                                </Typography>
                            </div>
                            <div>
                                <Button onClick={() => window.open(imageMetadata.externalUrl, '_blank')}>
                                    {intl.formatMessage({id: INTL_DATA.MARKETPLACE})}
                                </Button>
                            </div>
                        </div>
                        <div className={styles.description}>
                            <Typography>
                                {imageMetadata.description}
                            </Typography>
                        </div>
                        <div className={styles.certificate}>
                            {
                                !imageCertificate.id ?
                                    <EmptyCertificate
                                        certificationPrice={certificationPrice}
                                        isFetching={isFetching}
                                        onRunCheck={runCheckImage}
                                    /> : !imageCertificate.date ?
                                        <ProgressCertificate certificate={imageCertificate}/> :
                                        <Certificate
                                            certificate={imageCertificate}
                                            isFetching={isFetching}
                                            onRunCheck={runCheckImage}
                                        />

                            }
                        </div>
                    </div>
                </div>
            </Container>
            <ErrorMessage
                isOpen={!!errorMessage}
                message={errorMessage}
                onClose={() => setErrorMessage('')}
            />
            <SuccessMessage
                isOpen={!!successMessage}
                message={successMessage}
                onClose={() => setSuccessMessage('')}
            />
        </>
    );
};

export default Check;
