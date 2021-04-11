export default function(date: Date | null) {
    return Object.prototype.toString.call(date) === '[object Date]';
}
