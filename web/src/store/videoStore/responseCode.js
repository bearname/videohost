const RESPONSE_CODES = {
  SuccessAddLike: 0,
  SuccessAddDislike: 1,
  SuccessDeleteLike: 2,
  SuccessDeleteDisLike: 3,
  ErrFailedDeleteLike: 102,
  ErrFailedAddLike: 103,
  ErrAlreadyLike: 104,
  ErrAlreadyDisLike: 105,
  ErrInternalServer: 500,
};

export default RESPONSE_CODES;
