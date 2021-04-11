export default function(date?: Date) {
    const res = date || new Date();

    res.setFullYear(res.getFullYear() - 1);

    return res;
}
