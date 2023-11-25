const API_PATH = '';

export const getAddSurvey = (key, value) => {
    return fetch(`${API_PATH}/survey/add`, {
        method: 'POST', headers: {
            'Content-Type': 'application/json;charset=utf-8',
        }, credentials: 'include',
        body: JSON.stringify({
            attribute: key,
            metric: value,
        }),
    }).then(response => {
        if (!response.ok) {
            throw new Error('Ответ сервера не 200');
        }

        return response.json();
    })
        .then(data => {
            return data;
        })
        .catch(error => {
            throw new Error(error);
        });
};