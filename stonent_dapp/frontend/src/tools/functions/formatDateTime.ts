export default function(date: any): string {
    date = new Date(date);

    if (isNaN(date.valueOf())) {
        return '';
    }

    const hours = date.getHours() < 10 ? '0' + date.getHours() : date.getHours();
    const minutes = date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes();
    const day = date.getDate() < 10 ? '0' + date.getDate() : date.getDate();
    const month = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1;
    const year = date.getFullYear();

    return `${hours}:${minutes} ${day}.${month}.${year}`;
}
