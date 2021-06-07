const videosUtil = {
    getElapsedString(uploadedDate) {
        let elapsed = (Date.now() - Date.parse(uploadedDate)) / 1000;

        elapsed = Math.round(elapsed);

        if (elapsed / 60 < 1) {
            return elapsed + " секунд назад"
        }
        elapsed = Math.round(elapsed / 60);
        if (elapsed / 60 < 1) {
            return elapsed + " минут назад"
        }
        elapsed = Math.round(elapsed / 24);
        if (elapsed / 24 < 1) {
            return elapsed + " часов назад"
        }
        let weeks = Math.round(elapsed / 7);
        if (weeks / 7 < 1) {
            return weeks + " недель назад"
        }
        let months = Math.round(elapsed / 30);
        if (months / 30 < 1) {
            return months + " месяцев назад"
        }
        let years = Math.round(elapsed / 365);
        if (years / 365 < 1) {
            return years + " лет назад"
        }
    },
    updateThumbnail(part, index) {
        this[index].thumbnail = process.env.VUE_APP_VIDEO_API + "/" + this[index].thumbnail
    },

};

export default videosUtil