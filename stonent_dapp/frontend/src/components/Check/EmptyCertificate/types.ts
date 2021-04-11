export interface IProps {
    certificationPrice: number | null;
    isFetching: boolean;

    onRunCheck(): void;
}
