import { useState } from 'react';
import { useContext } from 'react';
import { MangaContext } from '../../Context/MangaContext';
import DisplayImg from '../displayImg';
import NULL_CONSTANTS from '../../utilities/nullConstants';
import ReactPaginate from 'react-paginate';

import Link from 'next/link';
function ByAll() {
  const { mostViewedAll, server } =
    useContext(MangaContext);
  const [pageNumber, setPageNumber] = useState(0);
  const bookPerPage = 12;
  const pageVisited = pageNumber * bookPerPage;
  const displayDatas = mostViewedAll
    .slice(pageVisited, pageVisited + bookPerPage)
    .map((listObject, index) => (
      <Link
        href={`/mangas/${listObject.id}`}
        key={index}
        passHref
      >
        <div
          className="col-6 col-lg-3 col-md-4 col-xl-2"
          data-aos="fade-up"
        >
          <DisplayImg
            bgColor="green"
            srcImg={
              listObject.image
                ? `${server}/image/${listObject.image}`
                : NULL_CONSTANTS.BOOK_GROUP_IMAGE
            }
            text={listObject.views + ' lượt đọc'}
            title={listObject.title}
            height="282px"
          ></DisplayImg>
        </div>
      </Link>
    ));
  const pageCount = Math.ceil(
    mostViewedAll.length / bookPerPage
  );
  const changePage = ({ selected }) => {
    setPageNumber(selected);
  };

  return (
    <div>
      <div className="row">{displayDatas}</div>
      <ReactPaginate
        breakLabel="..."
        previousLabel="Trước"
        nextLabel="Sau"
        pageCount={pageCount}
        onPageChange={changePage}
        pageClassName="page-item"
        pageLinkClassName="page-link"
        previousClassName="page-item"
        previousLinkClassName="page-link"
        nextClassName="page-item"
        nextLinkClassName="page-link"
        breakClassName="page-item"
        breakLinkClassName="page-link"
        containerClassName="pagination"
        activeClassName="active"
        renderOnZeroPageCount={null}
      ></ReactPaginate>
    </div>
  );
}

export default ByAll;
